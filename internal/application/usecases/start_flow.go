package usecases

import (
	"context"
	"fmt"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// StartFlowUseCase inicia un flujo para una conversación
type StartFlowUseCase struct {
	flowEngine ports.FlowEngine
	flowRepo   ports.FlowRepository
	logger     ports.Logger
}

// NewStartFlowUseCase crea un nuevo caso de uso
func NewStartFlowUseCase(
	flowEngine ports.FlowEngine,
	flowRepo ports.FlowRepository,
	logger ports.Logger,
) *StartFlowUseCase {
	return &StartFlowUseCase{
		flowEngine: flowEngine,
		flowRepo:   flowRepo,
		logger:     logger,
	}
}

// StartFlowRequest es la solicitud para iniciar un flujo
type StartFlowRequest struct {
	ConversationID string
	FlowID         string // Si está vacío, usar flujo por defecto
	TenantID       string
	InstanceID     string
}

// Execute inicia el flujo
func (uc *StartFlowUseCase) Execute(ctx context.Context, req StartFlowRequest) (*entities.FlowSession, error) {
	uc.logger.Info(fmt.Sprintf("Starting flow for conversation: %s", req.ConversationID))

	flowID := req.FlowID

	// Si no se especifica flujo, buscar el flujo por defecto
	if flowID == "" {
		defaultFlow, err := uc.flowRepo.FindDefault(ctx, req.InstanceID)
		if err != nil {
			return nil, fmt.Errorf("no default flow found: %w", err)
		}
		flowID = defaultFlow.ID
		uc.logger.Info(fmt.Sprintf("Using default flow: %s", flowID))
	}

	// Iniciar flujo
	session, err := uc.flowEngine.StartFlow(ctx, req.ConversationID, flowID, req.TenantID, req.InstanceID)
	if err != nil {
		return nil, fmt.Errorf("error starting flow: %w", err)
	}

	uc.logger.Info(fmt.Sprintf("Flow started successfully, session ID: %s", session.ID))

	return session, nil
}


