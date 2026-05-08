package usecases

import (
	"context"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// GetConnectionStatusUseCaseImpl implementa el caso de uso de obtener estado de conexión
type GetConnectionStatusUseCaseImpl struct {
	messagingService ports.MessagingService
	logger           ports.Logger
}

// NewGetConnectionStatusUseCase crea una nueva instancia del caso de uso
func NewGetConnectionStatusUseCase(
	messagingService ports.MessagingService,
	logger ports.Logger,
) ports.GetConnectionStatusUseCase {
	return &GetConnectionStatusUseCaseImpl{
		messagingService: messagingService,
		logger:           logger,
	}
}

// Execute ejecuta el caso de uso de obtener estado de conexión
func (uc *GetConnectionStatusUseCaseImpl) Execute(ctx context.Context) (*entities.Connection, error) {
	connection, err := uc.messagingService.GetConnection(ctx)
	if err != nil {
		uc.logger.Error("Error al obtener estado de conexión", "error", err)
		return nil, err
	}

	uc.logger.Debug("Estado de conexión obtenido",
		"isConnected", connection.IsConnected,
		"isLoggedIn", connection.IsLoggedIn)

	return connection, nil
}
