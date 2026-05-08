package http

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"whatsapp-api-go/internal/application/usecases"
	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// WebhookHandler maneja los webhooks de WhatsApp Cloud API
type WebhookHandler struct {
	verifyToken               string
	appSecret                 string
	sendMessageUseCase        ports.SendMessageUseCase
	messageRepo               ports.MessageRepository
	instanceID                string // WABA_PHONE_ID
	logger                    ports.Logger
	startFlowUseCase          *usecases.StartFlowUseCase
	processFlowMessageUseCase *usecases.ProcessFlowMessageUseCase

	// Idempotencia simple (en producción usar Redis con TTL)
	seenMu      sync.RWMutex
	seenWamids  map[string]time.Time
	cleanupTick *time.Ticker
}

// NewWebhookHandler crea un nuevo handler de webhooks
func NewWebhookHandler(
	verifyToken, appSecret, instanceID string,
	sendMessageUseCase ports.SendMessageUseCase,
	messageRepo ports.MessageRepository,
	logger ports.Logger,
	startFlowUseCase *usecases.StartFlowUseCase,
	processFlowMessageUseCase *usecases.ProcessFlowMessageUseCase,
) *WebhookHandler {
	h := &WebhookHandler{
		verifyToken:               verifyToken,
		appSecret:                 appSecret,
		sendMessageUseCase:        sendMessageUseCase,
		messageRepo:               messageRepo,
		instanceID:                instanceID,
		logger:                    logger,
		startFlowUseCase:          startFlowUseCase,
		processFlowMessageUseCase: processFlowMessageUseCase,
		seenWamids:                make(map[string]time.Time),
		cleanupTick:               time.NewTicker(10 * time.Minute),
	}

	// Cleanup periódico de wamids viejos (> 1 hora)
	go h.cleanupOldWamids()

	return h
}

// VerifyWebhook maneja GET /webhook (validación de Meta)
func (h *WebhookHandler) VerifyWebhook(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	h.logger.Info("Webhook verification request", "mode", mode, "token_match", token == h.verifyToken)

	if mode == "subscribe" && token == h.verifyToken {
		h.logger.Info("✅ Webhook verified successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		return
	}

	h.logger.Warn("❌ Webhook verification failed")
	w.WriteHeader(http.StatusForbidden)
}

// ReceiveWebhook maneja POST /webhook (eventos de Meta)
func (h *WebhookHandler) ReceiveWebhook(w http.ResponseWriter, r *http.Request) {
	// Leer body completo para validar firma
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Error reading body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validar firma X-Hub-Signature-256
	signature := r.Header.Get("X-Hub-Signature-256")
	if !h.validateSignature(body, signature) {
		h.logger.Error("❌ Invalid signature")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// ACK rápido (Meta espera 200 inmediato)
	w.WriteHeader(http.StatusOK)

	// Procesar webhook en goroutine (no bloquear respuesta)
	go h.processWebhook(body)
}

// validateSignature valida la firma HMAC-SHA256
func (h *WebhookHandler) validateSignature(body []byte, signature string) bool {
	if h.appSecret == "" || signature == "" {
		return false
	}

	// Firma esperada: sha256=<hex>
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	expectedSig := signature[7:] // Quitar "sha256="

	// Calcular HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(h.appSecret))
	mac.Write(body)
	calculatedSig := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(calculatedSig), []byte(expectedSig))
}

// processWebhook procesa el webhook de forma asíncrona
func (h *WebhookHandler) processWebhook(body []byte) {
	var webhook struct {
		Object string `json:"object"`
		Entry  []struct {
			ID      string `json:"id"`
			Changes []struct {
				Value struct {
					MessagingProduct string `json:"messaging_product"`
					Metadata         struct {
						DisplayPhoneNumber string `json:"display_phone_number"`
						PhoneNumberID      string `json:"phone_number_id"`
					} `json:"metadata"`
					Contacts []struct {
						Profile struct {
							Name string `json:"name"`
						} `json:"profile"`
						WAID string `json:"wa_id"`
					} `json:"contacts"`
					Messages []struct {
						From      string `json:"from"`
						ID        string `json:"id"` // wamid
						Timestamp string `json:"timestamp"`
						Type      string `json:"type"`
						Text      struct {
							Body string `json:"body"`
						} `json:"text"`
						Audio struct {
							ID       string `json:"id"` // media_id para descargar
							MimeType string `json:"mime_type"`
							SHA256   string `json:"sha256"`
						} `json:"audio"`
						Image struct {
							ID       string `json:"id"`
							MimeType string `json:"mime_type"`
							SHA256   string `json:"sha256"`
							Caption  string `json:"caption"`
						} `json:"image"`
						Video struct {
							ID       string `json:"id"`
							MimeType string `json:"mime_type"`
							SHA256   string `json:"sha256"`
							Caption  string `json:"caption"`
						} `json:"video"`
						Document struct {
							ID       string `json:"id"`
							MimeType string `json:"mime_type"`
							SHA256   string `json:"sha256"`
							Filename string `json:"filename"`
							Caption  string `json:"caption"`
						} `json:"document"`
						Interactive struct {
							Type        string `json:"type"`
							ButtonReply struct {
								ID    string `json:"id"`
								Title string `json:"title"`
							} `json:"button_reply"`
							ListReply struct {
								ID          string `json:"id"`
								Title       string `json:"title"`
								Description string `json:"description"`
							} `json:"list_reply"`
						} `json:"interactive"`
					} `json:"messages"`
					Statuses []struct {
						ID        string `json:"id"`
						Status    string `json:"status"`
						Timestamp string `json:"timestamp"`
					} `json:"statuses"`
				} `json:"value"`
				Field string `json:"field"`
			} `json:"changes"`
		} `json:"entry"`
	}

	if err := json.Unmarshal(body, &webhook); err != nil {
		h.logger.Error("Error parsing webhook", "error", err)
		return
	}

	h.logger.Debug("Webhook received", "object", webhook.Object)

	// Procesar entries
	for _, entry := range webhook.Entry {
		for _, change := range entry.Changes {
			// Solo procesar mensajes (ignorar statuses para evitar loops)
			messages := change.Value.Messages

			metadata := change.Value.Metadata

			for _, msg := range messages {
				// Idempotencia: verificar dedup_key en MongoDB
				dedupKey := h.instanceID + "|" + msg.ID
				exists, _ := h.messageRepo.ExistsByDedupKey(context.Background(), dedupKey)
				if exists {
					h.logger.Debug("Mensaje ya procesado (duplicado)", "wamid", msg.ID)
					continue
				}

				// También verificar memoria (fallback)
				if h.isWamidSeen(msg.ID) {
					h.logger.Debug("Wamid ya en memoria", "wamid", msg.ID)
					continue
				}

				// Marcar como visto
				h.markWamidSeen(msg.ID)

				// Log del mensaje
				h.logger.Info("📨 Mensaje entrante",
					"from", msg.From,
					"wamid", msg.ID,
					"type", msg.Type)

				// Crear mensaje entrante según el tipo
				var incomingMsg *entities.Message

				switch msg.Type {
				case "text":
					incomingMsg = entities.NewIncomingMessage(
						h.instanceID,
						msg.ID,
						msg.From,
						metadata.PhoneNumberID,
						msg.Text.Body,
					)

				case "audio":
					// Crear MediaContent con el media_id de WhatsApp
					media := &entities.MediaContent{
						MimeType: msg.Audio.MimeType,
						SHA256:   msg.Audio.SHA256,
						// NO guardar Data aún, solo la referencia al media_id de WhatsApp
						Storage: &entities.MediaStorage{
							Provider: "whatsapp",
							Key:      msg.Audio.ID, // Este es el media_id para descargar después
						},
					}
					incomingMsg = entities.NewIncomingMediaMessage(
						h.instanceID,
						msg.ID,
						msg.From,
						metadata.PhoneNumberID,
						"audio",
						media,
						nil,
					)

				case "image":
					media := &entities.MediaContent{
						MimeType: msg.Image.MimeType,
						SHA256:   msg.Image.SHA256,
						Caption:  msg.Image.Caption,
						Storage: &entities.MediaStorage{
							Provider: "whatsapp",
							Key:      msg.Image.ID,
						},
					}
					incomingMsg = entities.NewIncomingMediaMessage(
						h.instanceID,
						msg.ID,
						msg.From,
						metadata.PhoneNumberID,
						"image",
						media,
						nil,
					)

				case "video":
					media := &entities.MediaContent{
						MimeType: msg.Video.MimeType,
						SHA256:   msg.Video.SHA256,
						Caption:  msg.Video.Caption,
						Storage: &entities.MediaStorage{
							Provider: "whatsapp",
							Key:      msg.Video.ID,
						},
					}
					incomingMsg = entities.NewIncomingMediaMessage(
						h.instanceID,
						msg.ID,
						msg.From,
						metadata.PhoneNumberID,
						"video",
						media,
						nil,
					)

				case "document":
					media := &entities.MediaContent{
						MimeType: msg.Document.MimeType,
						SHA256:   msg.Document.SHA256,
						FileName: msg.Document.Filename,
						Caption:  msg.Document.Caption,
						Storage: &entities.MediaStorage{
							Provider: "whatsapp",
							Key:      msg.Document.ID,
						},
					}
					incomingMsg = entities.NewIncomingMediaMessage(
						h.instanceID,
						msg.ID,
						msg.From,
						metadata.PhoneNumberID,
						"document",
						media,
						nil,
					)

				case "interactive":
					// Botones o listas
					interactive := make(map[string]interface{})
					if msg.Interactive.Type == "button_reply" && msg.Interactive.ButtonReply.ID != "" {
						interactive["button_reply"] = map[string]interface{}{
							"id":    msg.Interactive.ButtonReply.ID,
							"title": msg.Interactive.ButtonReply.Title,
						}
					} else if msg.Interactive.Type == "list_reply" && msg.Interactive.ListReply.ID != "" {
						interactive["list_reply"] = map[string]interface{}{
							"id":          msg.Interactive.ListReply.ID,
							"title":       msg.Interactive.ListReply.Title,
							"description": msg.Interactive.ListReply.Description,
						}
					}
					incomingMsg = entities.NewIncomingMediaMessage(
						h.instanceID,
						msg.ID,
						msg.From,
						metadata.PhoneNumberID,
						"interactive",
						nil,
						interactive,
					)

				default:
					h.logger.Warn("Tipo de mensaje no soportado", "type", msg.Type)
					continue
				}

				// Agregar raw_min para trazabilidad
				incomingMsg.RawMin = &entities.RawMinimal{
					EntryID:     entry.ID,
					ChangeField: change.Field,
					Metadata: map[string]interface{}{
						"display_phone_number": metadata.DisplayPhoneNumber,
					},
				}

				// Guardar mensaje entrante en MongoDB
				// WORKAROUND: No guardar mensajes de audio temporalmente (problema de serialización BSON)
				if msg.Type != "audio" {
					h.logger.Info("💾 Intentando guardar mensaje", "wamid", msg.ID, "type", msg.Type, "has_media", incomingMsg.MessageData.Media != nil)
					if err := h.messageRepo.Save(context.Background(), incomingMsg); err != nil {
						h.logger.Error("❌ Error guardando mensaje entrante", "error", err, "wamid", msg.ID, "error_type", fmt.Sprintf("%T", err))
					} else {
						h.logger.Info("✅ Mensaje guardado en MongoDB", "wamid", msg.ID)
					}
				} else {
					h.logger.Info("⏭️ Skipping save for audio message (will process directly)", "wamid", msg.ID)
				}

				// INTEGRACIÓN DE FLUJOS
				// Intentar procesar mensaje en flujo activo
				if h.processFlowMessageUseCase != nil {
					err := h.processFlowMessageUseCase.Execute(context.Background(), incomingMsg)
					if err != nil {
						h.logger.Error("Error procesando mensaje en flujo", "error", err, "wamid", msg.ID)
					} else {
						h.logger.Info("✅ Mensaje procesado en flujo", "wamid", msg.ID)
						continue // Saltar respuesta automática si se procesó en flujo
					}
				}

				// Si no se procesó en flujo, iniciar flujo por defecto
				if h.startFlowUseCase != nil {
					h.logger.Info("No hay sesión activa, iniciando flujo por defecto", "from", msg.From)
					_, err := h.startFlowUseCase.Execute(context.Background(), usecases.StartFlowRequest{
						ConversationID: incomingMsg.ConversationID, // FIX: usar ConversationID completo (phone@instance)
						FlowID:         "",                         // Usar flujo por defecto
						TenantID:       "default",
						InstanceID:     h.instanceID,
					})
					if err != nil {
						h.logger.Error("Error iniciando flujo", "error", err, "from", msg.From)

						// Fallback: respuesta automática simple
						if msg.Type == "text" && msg.Text.Body != "" {
							responseText := fmt.Sprintf("✅ Recibido: %s", msg.Text.Body)
							_, _ = h.sendMessageUseCase.Execute(context.Background(), msg.From, responseText)
						}
					}
				} else {
					// Sin flujos configurados, usar respuesta automática
					if msg.Type == "text" && msg.Text.Body != "" {
						responseText := fmt.Sprintf("✅ Recibido: %s", msg.Text.Body)
						_, err := h.sendMessageUseCase.Execute(context.Background(), msg.From, responseText)
						if err != nil {
							h.logger.Error("Error enviando respuesta", "error", err, "to", msg.From)
						}
					}
				}
			}

			// Procesar statuses (actualizaciones de estado de mensajes salientes)
			for _, status := range change.Value.Statuses {
				h.logger.Debug("📊 Status update",
					"wamid", status.ID,
					"status", status.Status,
					"timestamp", status.Timestamp)

				// Buscar mensaje por wamid
				message, err := h.messageRepo.FindByID(context.Background(), status.ID)
				if err == nil && message != nil {
					// Actualizar estado
					message.UpdateStatus(status.Status, status.ID)

					// Guardar actualización
					if err := h.messageRepo.Save(context.Background(), message); err != nil {
						h.logger.Error("Error actualizando status", "error", err, "wamid", status.ID)
					} else {
						h.logger.Info("✅ Status actualizado", "wamid", status.ID, "status", status.Status)
					}
				}
			}
		}
	}
}

// isWamidSeen verifica si ya vimos este wamid
func (h *WebhookHandler) isWamidSeen(wamid string) bool {
	h.seenMu.RLock()
	defer h.seenMu.RUnlock()
	_, seen := h.seenWamids[wamid]
	return seen
}

// markWamidSeen marca un wamid como visto
func (h *WebhookHandler) markWamidSeen(wamid string) {
	h.seenMu.Lock()
	defer h.seenMu.Unlock()
	h.seenWamids[wamid] = time.Now()
}

// cleanupOldWamids limpia wamids > 1 hora
func (h *WebhookHandler) cleanupOldWamids() {
	for range h.cleanupTick.C {
		h.seenMu.Lock()
		now := time.Now()
		for wamid, seenAt := range h.seenWamids {
			if now.Sub(seenAt) > 1*time.Hour {
				delete(h.seenWamids, wamid)
			}
		}
		h.seenMu.Unlock()
	}
}
