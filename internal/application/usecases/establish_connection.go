package usecases

import (
	"context"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// EstablishConnectionUseCaseImpl implementa el caso de uso de establecer conexión
type EstablishConnectionUseCaseImpl struct {
	messagingService ports.MessagingService
	sessionRepo      ports.SessionRepository
	logger           ports.Logger
}

// NewEstablishConnectionUseCase crea una nueva instancia del caso de uso
func NewEstablishConnectionUseCase(
	messagingService ports.MessagingService,
	sessionRepo ports.SessionRepository,
	logger ports.Logger,
) ports.EstablishConnectionUseCase {
	return &EstablishConnectionUseCaseImpl{
		messagingService: messagingService,
		sessionRepo:      sessionRepo,
		logger:           logger,
	}
}

// Execute ejecuta el caso de uso de establecer conexión
func (uc *EstablishConnectionUseCaseImpl) Execute(ctx context.Context) (*entities.Connection, error) {
	uc.logger.Info("Estableciendo conexión con WhatsApp...")

	// Conectar al servicio de mensajería
	connection, err := uc.messagingService.Connect(ctx)
	if err != nil {
		uc.logger.Error("Error al conectar con WhatsApp", "error", err)
		return nil, err
	}

	// Si hay sesión activa, guardarla
	if connection.Session != nil {
		if err := uc.sessionRepo.Save(ctx, connection.Session); err != nil {
			uc.logger.Warn("Error al guardar sesión", "error", err)
			// No retornamos error porque la conexión ya está establecida
		}
	}

	uc.logger.Info("Conexión establecida exitosamente")
	return connection, nil
}
