package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"sync"
	"time"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// CloudAPIAdapter implementa MessagingService usando WhatsApp Business Cloud API
type CloudAPIAdapter struct {
	phoneNumberID string
	accessToken   string
	apiVersion    string
	baseURL       string
	httpClient    *http.Client
	logger        ports.Logger

	// Rate limiting por usuario (pair rate limit: 1 msg / 6s)
	rateLimitMu sync.Mutex
	lastSent    map[string]time.Time // phone -> último envío
	rateLimit   time.Duration        // 6 segundos
}

// NewCloudAPIAdapter crea un nuevo adaptador de Cloud API
func NewCloudAPIAdapter(phoneNumberID, accessToken, apiVersion string, logger ports.Logger) *CloudAPIAdapter {
	return &CloudAPIAdapter{
		phoneNumberID: phoneNumberID,
		accessToken:   accessToken,
		apiVersion:    apiVersion,
		baseURL:       fmt.Sprintf("https://graph.facebook.com/%s", apiVersion),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger:    logger,
		lastSent:  make(map[string]time.Time),
		rateLimit: 6 * time.Second, // Meta recomienda 1 mensaje cada 6 segundos por usuario
	}
}

// SendMessage envía un mensaje usando Cloud API
func (a *CloudAPIAdapter) SendMessage(ctx context.Context, message *entities.Message) error {
	// Pair rate limiting
	if err := a.checkRateLimit(message.To); err != nil {
		message.SetError(429, "Rate limit hit", "Pair rate limit: 1 mensaje cada 6s")
		return err
	}

	url := fmt.Sprintf("%s/%s/messages", a.baseURL, a.phoneNumberID)

	// Construir payload según tipo de mensaje
	payload, err := a.buildPayload(message)
	if err != nil {
		return fmt.Errorf("error al construir payload: %w", err)
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al serializar mensaje: %w", err)
	}

	// LOG: Ver payload exacto que se envía
	a.logger.Info("📤 Enviando a WhatsApp Cloud API", "type", message.MessageData.Type, "payload", string(jsonData))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error al crear request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.accessToken))

	resp, err := a.httpClient.Do(req)
	if err != nil {
		message.SetError(500, "HTTP error", err.Error())
		return fmt.Errorf("error al enviar mensaje: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		a.logger.Error("Error de Cloud API", "status", resp.StatusCode, "body", string(body))

		// Parsear error de Meta
		var metaError struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &metaError); err == nil {
			message.SetError(metaError.Error.Code, "Meta API Error", metaError.Error.Message)
		} else {
			message.SetError(resp.StatusCode, "Cloud API error", string(body))
		}

		return fmt.Errorf("cloud API error: %s (status %d)", string(body), resp.StatusCode)
	}

	// Parsear respuesta para obtener wamid
	var result struct {
		Messages []struct {
			ID string `json:"id"` // wamid
		} `json:"messages"`
	}

	if err := json.Unmarshal(body, &result); err == nil && len(result.Messages) > 0 {
		wamid := result.Messages[0].ID
		message.ID = wamid
		message.MessageData.ID = wamid
		message.DedupKey = a.phoneNumberID + "|" + wamid
		message.UpdateStatus("sent", wamid)

		a.logger.Info("Mensaje enviado", "wamid", wamid, "to", message.To)
	}

	// Actualizar último envío para rate limiting
	a.updateLastSent(message.To)

	return nil
}

// buildPayload construye el payload según el tipo de mensaje
func (a *CloudAPIAdapter) buildPayload(message *entities.Message) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                message.To,
	}

	// Tipo de mensaje
	msgType := message.MessageData.Type
	payload["type"] = msgType

	switch msgType {
	case "text":
		if message.MessageData.Text != nil {
			payload["text"] = map[string]string{
				"body": message.MessageData.Text.Body,
			}
		}

	case "image", "video", "audio", "document":
		if message.MessageData.Media != nil {
			mediaPayload := map[string]interface{}{}

			// Si tiene datos binarios, primero subir a WhatsApp
			if len(message.MessageData.Media.Data) > 0 {
				mediaID, err := a.uploadMedia(message.MessageData.Media.Data, message.MessageData.Media.MimeType)
				if err != nil {
					a.logger.Error("Error uploading media to WhatsApp", "error", err)
					return nil, fmt.Errorf("error uploading media: %w", err)
				}
				mediaPayload["id"] = mediaID
			} else if message.MessageData.Media.Storage != nil {
				// Si tiene URL pública, usar link
				mediaPayload["link"] = message.MessageData.Media.Storage.PublicURL
			}

			if message.MessageData.Media.Caption != "" {
				mediaPayload["caption"] = message.MessageData.Media.Caption
			}
			payload[msgType] = mediaPayload
		}

	case "location":
		if message.MessageData.Location != nil {
			payload["location"] = map[string]interface{}{
				"latitude":  message.MessageData.Location.Latitude,
				"longitude": message.MessageData.Location.Longitude,
				"name":      message.MessageData.Location.Name,
				"address":   message.MessageData.Location.Address,
			}
		}

	case "interactive":
		// Botones interactivos, listas, etc.
		if message.MessageData.Interactive != nil {
			payload["interactive"] = message.MessageData.Interactive
		}
	}

	// Context (reply_to)
	if message.MessageData.Context != nil {
		payload["context"] = map[string]string{
			"message_id": message.MessageData.Context.MessageID,
		}
	}

	return payload, nil
}

// uploadMedia sube un archivo de medios a WhatsApp Cloud API usando multipart/form-data
func (a *CloudAPIAdapter) uploadMedia(data []byte, mimeType string) (string, error) {
	url := fmt.Sprintf("%s/%s/media", a.baseURL, a.phoneNumberID)

	// Crear buffer para multipart
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Agregar messaging_product
	writer.WriteField("messaging_product", "whatsapp")

	// Determinar extensión y mime type correcto
	extension := ".ogg"
	fileMimeType := "audio/ogg; codecs=opus"

	if strings.Contains(mimeType, "audio/webm") {
		extension = ".ogg"
		fileMimeType = "audio/ogg; codecs=opus" // WhatsApp prefiere OGG Opus
	} else if strings.Contains(mimeType, "audio/ogg") || strings.Contains(mimeType, "audio/opus") {
		extension = ".ogg"
		fileMimeType = "audio/ogg; codecs=opus"
	} else if strings.Contains(mimeType, "audio/mpeg") || strings.Contains(mimeType, "audio/mp3") {
		extension = ".mp3"
		fileMimeType = "audio/mpeg"
	}

	// Crear parte del archivo con Content-Type específico
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="audio%s"`, extension))
	h.Set("Content-Type", fileMimeType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return "", fmt.Errorf("error creating form file: %w", err)
	}

	if _, err := part.Write(data); err != nil {
		return "", fmt.Errorf("error writing file data: %w", err)
	}

	// Agregar type
	writer.WriteField("type", mimeType)

	// Cerrar writer
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("error closing multipart writer: %w", err)
	}

	// Crear request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", fmt.Errorf("error creating upload request: %w", err)
	}

	// Headers
	req.Header.Set("Authorization", "Bearer "+a.accessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Ejecutar request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error uploading media: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading upload response: %w", err)
	}

	// Verificar status
	if resp.StatusCode != http.StatusOK {
		a.logger.Error("WhatsApp media upload failed", "status", resp.StatusCode, "body", string(body))
		return "", fmt.Errorf("media upload failed (status %d): %s", resp.StatusCode, string(body))
	}

	// Parsear respuesta
	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error parsing upload response: %w", err)
	}

	a.logger.Info("Media uploaded successfully", "media_id", result.ID)
	return result.ID, nil
}

// checkRateLimit verifica el pair rate limit (1 mensaje cada 6s por usuario)
func (a *CloudAPIAdapter) checkRateLimit(phone string) error {
	a.rateLimitMu.Lock()
	defer a.rateLimitMu.Unlock()

	if lastTime, exists := a.lastSent[phone]; exists {
		elapsed := time.Since(lastTime)
		if elapsed < a.rateLimit {
			waitTime := a.rateLimit - elapsed
			a.logger.Warn("Rate limit activo", "phone", phone, "wait", waitTime.Seconds())
			time.Sleep(waitTime)
		}
	}

	return nil
}

// updateLastSent actualiza el timestamp del último envío
func (a *CloudAPIAdapter) updateLastSent(phone string) {
	a.rateLimitMu.Lock()
	defer a.rateLimitMu.Unlock()
	a.lastSent[phone] = time.Now()
}

// GetMessageStatus obtiene el estado de un mensaje (usando webhooks)
func (a *CloudAPIAdapter) GetMessageStatus(ctx context.Context, messageID string) (string, error) {
	// En Cloud API, los estados llegan por webhooks (statuses[])
	// Esta función sería para consultar si necesitamos
	return "sent", nil
}

// Connect no es necesario en Cloud API (sin sesión persistente)
func (a *CloudAPIAdapter) Connect(ctx context.Context) (*entities.Connection, error) {
	connection := entities.NewConnection()
	connection.IsConnected = true
	connection.IsLoggedIn = true

	a.logger.Info("✅ Cloud API configurada correctamente")
	a.logger.Info("📱 Phone Number ID:", "id", a.phoneNumberID)
	a.logger.Info("🔗 API Version:", "version", a.apiVersion)

	return connection, nil
}

// Disconnect no hace nada en Cloud API
func (a *CloudAPIAdapter) Disconnect(ctx context.Context) error {
	a.logger.Info("Cloud API - no requiere desconexión")
	return nil
}

// GetConnection siempre retorna conectado en Cloud API
func (a *CloudAPIAdapter) GetConnection(ctx context.Context) (*entities.Connection, error) {
	connection := entities.NewConnection()
	connection.IsConnected = true
	connection.IsLoggedIn = true
	return connection, nil
}

// IsConnected siempre retorna true en Cloud API (sin sesión local)
func (a *CloudAPIAdapter) IsConnected(ctx context.Context) bool {
	return true
}
