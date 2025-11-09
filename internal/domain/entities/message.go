package entities

import "time"

// Message representa un mensaje de WhatsApp en el dominio
type Message struct {
	ID string `json:"_id" bson:"_id"` // MongoDB _id

	// Multi-tenant y routing
	TenantID   string `json:"tenant_id,omitempty" bson:"tenant_id,omitempty"`
	InstanceID string `json:"instance_id" bson:"instance_id"` // WABA_PHONE_ID
	Channel    string `json:"channel" bson:"channel"`         // whatsapp
	Provider   string `json:"provider" bson:"provider"`       // meta

	// Dirección y conversación
	Direction      string       `json:"direction" bson:"direction"`                                 // in | out
	ConversationID string       `json:"conversation_id" bson:"conversation_id"`                     // phone@instance
	WAConversation Conversation `json:"wa_conversation,omitempty" bson:"wa_conversation,omitempty"` // Meta 24h window

	// Participantes
	From      string `json:"from" bson:"from"`                                 // E.164 cliente
	To        string `json:"to" bson:"to"`                                     // E.164 WABA
	ContactID string `json:"contact_id,omitempty" bson:"contact_id,omitempty"` // referencia opcional

	// Contenido del mensaje
	MessageData MessageData `json:"message" bson:"message"`

	// Estado y tracking
	Status        string          `json:"status" bson:"status"` // received|queued|sent|delivered|read|failed
	StatusHistory []StatusHistory `json:"status_history,omitempty" bson:"status_history,omitempty"`
	Error         *MessageError   `json:"error,omitempty" bson:"error,omitempty"`

	// Flow engine (opcional)
	FlowState *FlowState `json:"flow_state,omitempty" bson:"flow_state,omitempty"`

	// Deduplicación e idempotencia
	DedupKey string `json:"dedup_key" bson:"dedup_key"` // instance_id|wamid

	// Raw data mínimo (trazabilidad)
	RawMin *RawMinimal `json:"raw_min,omitempty" bson:"raw_min,omitempty"`

	// Timestamps
	Timestamps MessageTimestamps `json:"timestamps" bson:"timestamps"`
}

// Conversation representa la ventana de conversación de 24h de Meta
type Conversation struct {
	ID        string     `json:"id,omitempty" bson:"id,omitempty"`                 // value.statuses[*].conversation.id
	Category  string     `json:"category,omitempty" bson:"category,omitempty"`     // marketing|utility|authentication|service
	Origin    string     `json:"origin,omitempty" bson:"origin,omitempty"`         // user_initiated|business_initiated
	ExpiresAt *time.Time `json:"expires_at,omitempty" bson:"expires_at,omitempty"` // 24h desde inicio
}

// MessageData contiene el contenido del mensaje
type MessageData struct {
	ID   string `json:"id" bson:"id"`     // wamid
	Type string `json:"type" bson:"type"` // text|image|video|audio|document|location|interactive|contacts|sticker

	// Texto
	Text *TextContent `json:"text,omitempty" bson:"text,omitempty"`

	// Interactive (botones, listas, flows)
	Interactive *InteractiveContent `json:"interactive,omitempty" bson:"interactive,omitempty"`

	// Media (imagen, video, audio, documento, sticker)
	Media *MediaContent `json:"media,omitempty" bson:"media,omitempty"`

	// Ubicación
	Location *LocationContent `json:"location,omitempty" bson:"location,omitempty"`

	// Contexto (reply_to)
	Context *MessageContext `json:"context,omitempty" bson:"context,omitempty"`
}

// TextContent representa contenido de texto
type TextContent struct {
	Body string `json:"body" bson:"body"`
}

// InteractiveContent representa contenido interactivo
type InteractiveContent struct {
	Type string `json:"type" bson:"type"` // button_reply|list_reply|nfm_reply

	ButtonReply *ButtonReply `json:"button_reply,omitempty" bson:"button_reply,omitempty"`
	ListReply   *ListReply   `json:"list_reply,omitempty" bson:"list_reply,omitempty"`
	NFMReply    *NFMReply    `json:"nfm_reply,omitempty" bson:"nfm_reply,omitempty"`
}

// ButtonReply representa respuesta a botón
type ButtonReply struct {
	ID    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
}

// ListReply representa respuesta a lista
type ListReply struct {
	ID          string `json:"id" bson:"id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

// NFMReply representa respuesta a Flow
type NFMReply struct {
	Name         string                 `json:"name" bson:"name"`
	Body         string                 `json:"body,omitempty" bson:"body,omitempty"`
	ResponseJSON map[string]interface{} `json:"response_json,omitempty" bson:"response_json,omitempty"`
}

// MediaContent representa contenido multimedia
type MediaContent struct {
	MimeType string `json:"mime_type" bson:"mime_type"`
	FileName string `json:"file_name,omitempty" bson:"file_name,omitempty"`
	SHA256   string `json:"sha256,omitempty" bson:"sha256,omitempty"`
	Size     int64  `json:"size,omitempty" bson:"size,omitempty"`
	Caption  string `json:"caption,omitempty" bson:"caption,omitempty"`

	// NO guardar binario en Mongo - usar storage externo
	Storage *MediaStorage `json:"storage,omitempty" bson:"storage,omitempty"`
}

// MediaStorage representa almacenamiento externo de media
type MediaStorage struct {
	Provider  string `json:"provider" bson:"provider"` // s3|gcs|azure
	Bucket    string `json:"bucket" bson:"bucket"`
	Key       string `json:"key" bson:"key"`
	PublicURL string `json:"public_url" bson:"public_url"`
}

// LocationContent representa ubicación
type LocationContent struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
	Name      string  `json:"name,omitempty" bson:"name,omitempty"`
	Address   string  `json:"address,omitempty" bson:"address,omitempty"`
}

// MessageContext representa contexto (reply_to)
type MessageContext struct {
	MessageID string `json:"message_id" bson:"message_id"` // wamid del mensaje original
	From      string `json:"from,omitempty" bson:"from,omitempty"`
}

// StatusHistory representa historial de estados
type StatusHistory struct {
	Status     string    `json:"status" bson:"status"`
	Timestamp  time.Time `json:"ts" bson:"ts"`
	ProviderID string    `json:"provider_id,omitempty" bson:"provider_id,omitempty"` // wamid si aplica
}

// MessageError representa error en mensaje
type MessageError struct {
	Code    int    `json:"code" bson:"code"`
	Title   string `json:"title" bson:"title"`
	Details string `json:"details,omitempty" bson:"details,omitempty"`
}

// FlowState representa estado de flujo conversacional
type FlowState struct {
	FlowID  string                 `json:"flow_id" bson:"flow_id"`
	Version int                    `json:"version,omitempty" bson:"version,omitempty"`
	Step    string                 `json:"step,omitempty" bson:"step,omitempty"`
	Context map[string]interface{} `json:"context,omitempty" bson:"context,omitempty"`
}

// RawMinimal representa datos raw mínimos para trazabilidad
type RawMinimal struct {
	EntryID     string                 `json:"entry_id,omitempty" bson:"entry_id,omitempty"`
	ChangeField string                 `json:"change_field,omitempty" bson:"change_field,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// MessageTimestamps representa timestamps del mensaje
type MessageTimestamps struct {
	CreatedAt   time.Time  `json:"created_at" bson:"created_at"`                       // cuando se guardó
	ReceivedAt  *time.Time `json:"received_at,omitempty" bson:"received_at,omitempty"` // in
	QueuedAt    *time.Time `json:"queued_at,omitempty" bson:"queued_at,omitempty"`     // out
	SentAt      *time.Time `json:"sent_at,omitempty" bson:"sent_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty" bson:"delivered_at,omitempty"`
	ReadAt      *time.Time `json:"read_at,omitempty" bson:"read_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at" bson:"updated_at"`
}

// NewIncomingMessage crea un nuevo mensaje entrante
func NewIncomingMessage(instanceID, wamid, from, to, body string) *Message {
	now := time.Now()
	conversationID := from + "@" + instanceID

	return &Message{
		ID:             wamid,
		InstanceID:     instanceID,
		Channel:        "whatsapp",
		Provider:       "meta",
		Direction:      "in",
		ConversationID: conversationID,
		From:           from,
		To:             to,
		MessageData: MessageData{
			ID:   wamid,
			Type: "text",
			Text: &TextContent{Body: body},
		},
		Status:   "received",
		DedupKey: instanceID + "|" + wamid,
		StatusHistory: []StatusHistory{
			{Status: "received", Timestamp: now},
		},
		Timestamps: MessageTimestamps{
			CreatedAt:  now,
			ReceivedAt: &now,
			UpdatedAt:  now,
		},
	}
}

// NewOutgoingMessage crea un nuevo mensaje saliente
func NewOutgoingMessage(instanceID, to, body string) *Message {
	now := time.Now()
	conversationID := to + "@" + instanceID

	return &Message{
		InstanceID:     instanceID,
		Channel:        "whatsapp",
		Provider:       "meta",
		Direction:      "out",
		ConversationID: conversationID,
		From:           instanceID,
		To:             to,
		MessageData: MessageData{
			Type: "text",
			Text: &TextContent{Body: body},
		},
		Status: "queued",
		StatusHistory: []StatusHistory{
			{Status: "queued", Timestamp: now},
		},
		Timestamps: MessageTimestamps{
			CreatedAt: now,
			QueuedAt:  &now,
			UpdatedAt: now,
		},
	}
}

// UpdateStatus actualiza el estado del mensaje
func (m *Message) UpdateStatus(newStatus string, providerID ...string) {
	now := time.Now()
	m.Status = newStatus
	m.Timestamps.UpdatedAt = now

	// Agregar al historial
	history := StatusHistory{
		Status:    newStatus,
		Timestamp: now,
	}
	if len(providerID) > 0 {
		history.ProviderID = providerID[0]
	}
	m.StatusHistory = append(m.StatusHistory, history)

	// Actualizar timestamps específicos
	switch newStatus {
	case "sent":
		m.Timestamps.SentAt = &now
	case "delivered":
		m.Timestamps.DeliveredAt = &now
	case "read":
		m.Timestamps.ReadAt = &now
	}
}

// SetError establece un error en el mensaje
func (m *Message) SetError(code int, title, details string) {
	m.Status = "failed"
	m.Error = &MessageError{
		Code:    code,
		Title:   title,
		Details: details,
	}
	m.Timestamps.UpdatedAt = time.Now()
}

// Validate valida que el mensaje tenga los datos necesarios
func (m *Message) Validate() error {
	if m.To == "" {
		return ErrInvalidRecipient
	}
	if m.MessageData.Text == nil || m.MessageData.Text.Body == "" {
		if m.MessageData.Media == nil && m.MessageData.Location == nil {
			return ErrEmptyMessage
		}
	}
	return nil
}
