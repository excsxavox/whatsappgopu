package usecases

import (
	"context"
	"fmt"
	
	"whatsapp-api-go/internal/domain/ports"
)

// HandleWebhookUseCaseImpl implementa el caso de uso de procesar webhooks
type HandleWebhookUseCaseImpl struct {
	sendMessageUseCase ports.SendMessageUseCase
	logger             ports.Logger
}

// NewHandleWebhookUseCase crea una nueva instancia del caso de uso
func NewHandleWebhookUseCase(
	sendMessageUseCase ports.SendMessageUseCase,
	logger ports.Logger,
) ports.HandleWebhookUseCase {
	return &HandleWebhookUseCaseImpl{
		sendMessageUseCase: sendMessageUseCase,
		logger:             logger,
	}
}

// Execute ejecuta el caso de uso de procesar webhook
func (uc *HandleWebhookUseCaseImpl) Execute(ctx context.Context, payload map[string]interface{}) error {
	uc.logger.Info("Procesando webhook", "payload", payload)

	// Verificar si es una acción de enviar mensaje
	action, ok := payload["action"].(string)
	if !ok || action != "send_message" {
		uc.logger.Debug("Webhook no contiene acción send_message")
		return nil
	}

	// Extraer teléfono y mensaje
	phone, phoneOk := payload["phone"].(string)
	message, messageOk := payload["message"].(string)

	if !phoneOk || !messageOk || phone == "" || message == "" {
		err := fmt.Errorf("webhook inválido: faltan campos phone o message")
		uc.logger.Error("Webhook inválido", "error", err)
		return err
	}

	// Delegar al caso de uso de enviar mensaje
	_, err := uc.sendMessageUseCase.Execute(ctx, phone, message)
	if err != nil {
		uc.logger.Error("Error al procesar webhook", "error", err)
		return err
	}

	uc.logger.Info("Webhook procesado exitosamente")
	return nil
}

