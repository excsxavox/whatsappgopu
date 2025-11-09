package http

import (
	"encoding/json"
	"net/http"

	"whatsapp-api-go/internal/domain/ports"
)

// HTTPAdapter adapta las peticiones HTTP a los casos de uso del dominio
type HTTPAdapter struct {
	sendMessageUseCase         ports.SendMessageUseCase
	getConnectionStatusUseCase ports.GetConnectionStatusUseCase
	handleWebhookUseCase       ports.HandleWebhookUseCase
	logger                     ports.Logger
}

// NewHTTPAdapter crea un nuevo adaptador HTTP
func NewHTTPAdapter(
	sendMessageUseCase ports.SendMessageUseCase,
	getConnectionStatusUseCase ports.GetConnectionStatusUseCase,
	handleWebhookUseCase ports.HandleWebhookUseCase,
	logger ports.Logger,
) *HTTPAdapter {
	return &HTTPAdapter{
		sendMessageUseCase:         sendMessageUseCase,
		getConnectionStatusUseCase: getConnectionStatusUseCase,
		handleWebhookUseCase:       handleWebhookUseCase,
		logger:                     logger,
	}
}

// HealthHandler verifica el estado del servidor
func (h *HTTPAdapter) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "whatsapp-api-core",
	})
}

// SendMessageHandler maneja el envío de mensajes
func (h *HTTPAdapter) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Phone   string `json:"phone"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	if req.Phone == "" || req.Message == "" {
		http.Error(w, "Phone y Message son requeridos", http.StatusBadRequest)
		return
	}

	// Ejecutar caso de uso
	message, err := h.sendMessageUseCase.Execute(r.Context(), req.Phone, req.Message)
	if err != nil {
		h.logger.Error("Error al enviar mensaje", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "success",
		"message":    "Mensaje enviado correctamente",
		"message_id": message.ID,
	})
}

// StatusHandler muestra el estado de la conexión
func (h *HTTPAdapter) StatusHandler(w http.ResponseWriter, r *http.Request) {
	connection, err := h.getConnectionStatusUseCase.Execute(r.Context())
	if err != nil {
		h.logger.Error("Error al obtener estado", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"connected": connection.IsConnected,
		"logged_in": connection.IsLoggedIn,
	})
}

// WebhookHandler lo maneja webhook.go (Meta Cloud API)
