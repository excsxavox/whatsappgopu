package ports

import (
	"context"
	"whatsapp-api-go/internal/domain/entities"
)

// MessagingService define el puerto para enviar/recibir mensajes
// Este es un puerto de SALIDA (driving port)
type MessagingService interface {
	// SendMessage envía un mensaje de WhatsApp
	SendMessage(ctx context.Context, message *entities.Message) error

	// GetMessageStatus obtiene el estado de un mensaje
	GetMessageStatus(ctx context.Context, messageID string) (string, error)

	// Connect establece la conexión con WhatsApp
	Connect(ctx context.Context) (*entities.Connection, error)

	// Disconnect cierra la conexión
	Disconnect(ctx context.Context) error

	// GetConnection obtiene el estado actual de la conexión
	GetConnection(ctx context.Context) (*entities.Connection, error)

	// IsConnected verifica si hay conexión activa
	IsConnected(ctx context.Context) bool
}

// MessageRepository define el puerto para persistencia de mensajes
// Este es un puerto de SALIDA (driving port)
type MessageRepository interface {
	// Save guarda un mensaje
	Save(ctx context.Context, message *entities.Message) error

	// FindByID busca un mensaje por ID
	FindByID(ctx context.Context, messageID string) (*entities.Message, error)

	// FindByRecipient busca mensajes por destinatario
	FindByRecipient(ctx context.Context, recipient string, limit int) ([]*entities.Message, error)

	// UpdateStatus actualiza el estado de un mensaje
	UpdateStatus(ctx context.Context, messageID string, status string) error

	// ExistsByDedupKey verifica si existe un mensaje con ese dedup_key
	ExistsByDedupKey(ctx context.Context, dedupKey string) (bool, error)
}
