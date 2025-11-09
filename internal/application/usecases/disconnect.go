package usecases

import (
	"context"
	
	"whatsapp-api-go/internal/domain/ports"
)

// DisconnectUseCaseImpl implementa el caso de uso de desconectar
type DisconnectUseCaseImpl struct {
	messagingService ports.MessagingService
	sessionRepo      ports.SessionRepository
	logger           ports.Logger
}

// NewDisconnectUseCase crea una nueva instancia del caso de uso
func NewDisconnectUseCase(
	messagingService ports.MessagingService,
	sessionRepo ports.SessionRepository,
	logger ports.Logger,
) ports.DisconnectUseCase {
	return &DisconnectUseCaseImpl{
		messagingService: messagingService,
		sessionRepo:      sessionRepo,
		logger:           logger,
	}
}

// Execute ejecuta el caso de uso de desconectar
func (uc *DisconnectUseCaseImpl) Execute(ctx context.Context) error {
	uc.logger.Info("Desconectando de WhatsApp...")

	// Marcar sesión como inactiva
	session, err := uc.sessionRepo.FindActive(ctx)
	if err == nil && session != nil {
		_ = uc.sessionRepo.MarkAsInactive(ctx, session.ID)
	}

	// Desconectar del servicio
	if err := uc.messagingService.Disconnect(ctx); err != nil {
		uc.logger.Error("Error al desconectar", "error", err)
		return err
	}

	uc.logger.Info("Desconectado exitosamente")
	return nil
}

