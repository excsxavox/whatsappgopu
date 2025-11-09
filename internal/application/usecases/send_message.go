package usecases

import (
	"context"
	"fmt"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// SendMessageUseCaseImpl implementa el caso de uso de enviar mensajes
type SendMessageUseCaseImpl struct {
	messagingService ports.MessagingService
	messageRepo      ports.MessageRepository
	logger           ports.Logger
}

// NewSendMessageUseCase crea una nueva instancia del caso de uso
func NewSendMessageUseCase(
	messagingService ports.MessagingService,
	messageRepo ports.MessageRepository,
	logger ports.Logger,
) ports.SendMessageUseCase {
	return &SendMessageUseCaseImpl{
		messagingService: messagingService,
		messageRepo:      messageRepo,
		logger:           logger,
	}
}

// Execute ejecuta el caso de uso de enviar mensaje
func (uc *SendMessageUseCaseImpl) Execute(ctx context.Context, to, messageContent string) (*entities.Message, error) {
	// Verificar conexión
	if !uc.messagingService.IsConnected(ctx) {
		uc.logger.Error("No hay conexión activa con WhatsApp")
		return nil, entities.ErrNotConnected
	}

	// Crear el mensaje saliente con nueva estructura
	// Nota: instanceID debería venir del contexto/config, por ahora usamos genérico
	message := entities.NewOutgoingMessage("", to, messageContent)

	// Validar el mensaje
	if err := message.Validate(); err != nil {
		uc.logger.Error("Mensaje inválido", "error", err)
		return nil, err
	}

	// Enviar el mensaje a través del servicio de mensajería
	if err := uc.messagingService.SendMessage(ctx, message); err != nil {
		uc.logger.Error("Error al enviar mensaje", "error", err, "to", to)
		// Intentar guardar aunque haya fallado (ya tiene error seteado en SetError)
		_ = uc.messageRepo.Save(ctx, message)
		return nil, fmt.Errorf("error al enviar mensaje: %w", err)
	}

	// El mensaje ya fue actualizado con wamid y status "sent" en el adapter

	// Guardar el mensaje en el repositorio
	if err := uc.messageRepo.Save(ctx, message); err != nil {
		uc.logger.Warn("Error al guardar mensaje", "error", err)
		// No retornamos error porque el mensaje ya se envió
	}

	uc.logger.Info("Mensaje enviado exitosamente", "to", to, "wamid", message.ID)
	return message, nil
}
