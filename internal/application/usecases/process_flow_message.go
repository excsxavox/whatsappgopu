package usecases

import (
	"context"
	"fmt"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// ProcessFlowMessageUseCase procesa un mensaje en el contexto de un flujo
type ProcessFlowMessageUseCase struct {
	flowEngine  ports.FlowEngine
	sessionRepo ports.FlowSessionRepository
	logger      ports.Logger
}

// NewProcessFlowMessageUseCase crea un nuevo caso de uso
func NewProcessFlowMessageUseCase(
	flowEngine ports.FlowEngine,
	sessionRepo ports.FlowSessionRepository,
	logger ports.Logger,
) *ProcessFlowMessageUseCase {
	return &ProcessFlowMessageUseCase{
		flowEngine:  flowEngine,
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

// Execute procesa el mensaje en el flujo
func (uc *ProcessFlowMessageUseCase) Execute(ctx context.Context, message *entities.Message) error {
	uc.logger.Info(fmt.Sprintf("Processing message in flow: %s", message.ConversationID))

	// Buscar sesión activa
	session, err := uc.sessionRepo.FindActiveByConversation(ctx, message.ConversationID)
	if err != nil {
		return fmt.Errorf("error finding active session: %w", err)
	}

	if session == nil {
		uc.logger.Info("No active session found for this conversation")
		return fmt.Errorf("no active session found")
	}

	// Procesar mensaje en el contexto de la sesión
	err = uc.flowEngine.ProcessMessage(ctx, session, message)
	if err != nil {
		return fmt.Errorf("error processing message in flow: %w", err)
	}

	uc.logger.Info("Message processed successfully in flow")

	return nil
}


