package ports

import (
	"context"
	"whatsapp-api-go/internal/domain/entities"
)

// SendMessageUseCase define el puerto de ENTRADA para enviar mensajes
// Este es un puerto de ENTRADA (driven port) - define lo que el dominio ofrece al exterior
type SendMessageUseCase interface {
	Execute(ctx context.Context, to, message string) (*entities.Message, error)
}

// GetConnectionStatusUseCase define el puerto de ENTRADA para obtener estado
type GetConnectionStatusUseCase interface {
	Execute(ctx context.Context) (*entities.Connection, error)
}

// HandleWebhookUseCase define el puerto de ENTRADA para procesar webhooks
type HandleWebhookUseCase interface {
	Execute(ctx context.Context, payload map[string]interface{}) error
}

// EstablishConnectionUseCase define el puerto de ENTRADA para establecer conexión
type EstablishConnectionUseCase interface {
	Execute(ctx context.Context) (*entities.Connection, error)
}

// DisconnectUseCase define el puerto de ENTRADA para desconectar
type DisconnectUseCase interface {
	Execute(ctx context.Context) error
}

