package flow

import (
	"context"
	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// NodeProcessor es la interfaz común para todos los procesadores de nodos
type NodeProcessor interface {
	Process(ctx context.Context, session *entities.FlowSession, node *entities.FlowNode) (*ProcessResult, error)
}

// ProcessResult es el resultado del procesamiento de un nodo
type ProcessResult struct {
	WaitingForResponse bool   // Si debe esperar respuesta del usuario
	WaitingForVariable string // Nombre de la variable que espera
	NextNodeID         string // ID del siguiente nodo (para CONDITION)
	StopFlow           bool   // Si debe detener el flujo
	ErrorMessage       string // Mensaje de error si algo falló
}

// TTSService interface para servicios de texto a voz
type TTSService interface {
	TextToSpeech(text string) ([]byte, error)
}

// AIContextService interface for AI context-based responses
type AIContextService interface {
	RespondWithContext(userQuestion string, contexto *entities.Contexto) (string, error)
	IsQuestion(userMessage string) bool
}

// NodeProcessorFactory crea procesadores según el tipo de nodo
type NodeProcessorFactory struct {
	messagingService    ports.MessagingService
	logger              ports.Logger
	variableReplacer    *VariableReplacer
	aiValidationService AIValidationService
	aiContextService    AIContextService
	contextoRepo       ports.ContextoRepository
	flowRepo            ports.FlowRepository
	ttsService          TTSService
}

// NewNodeProcessorFactory crea una nueva factory
func NewNodeProcessorFactory(
	messagingService ports.MessagingService,
	logger ports.Logger,
	aiValidationService AIValidationService,
	aiContextService AIContextService,
	contextoRepo ports.ContextoRepository,
	flowRepo ports.FlowRepository,
	ttsService TTSService,
) *NodeProcessorFactory {
	return &NodeProcessorFactory{
		messagingService:    messagingService,
		logger:              logger,
		variableReplacer:    NewVariableReplacer(),
		aiValidationService: aiValidationService,
		aiContextService:    aiContextService,
		contextoRepo:       contextoRepo,
		flowRepo:            flowRepo,
		ttsService:          ttsService,
	}
}

// GetProcessor retorna el procesador adecuado según el tipo de nodo
func (f *NodeProcessorFactory) GetProcessor(nodeType string) NodeProcessor {
	switch nodeType {
	case "TEXT":
		return NewTextNodeProcessor(f.messagingService, f.logger, f.variableReplacer)
	case "BUTTONS":
		return NewButtonsNodeProcessor(f.messagingService, f.logger, f.variableReplacer)
	case "HTTP":
		return NewHttpNodeProcessor(f.logger, f.variableReplacer)
	case "CONDITION":
		return NewConditionNodeProcessor(f.logger, f.variableReplacer)
	case "RESPONSE":
		return NewResponseNodeProcessor(f.logger, f.variableReplacer, f.aiValidationService, f.aiContextService, f.contextoRepo, f.flowRepo, f.messagingService)
	case "AUDIO":
		return NewAudioNodeProcessor(f.messagingService, f.logger, f.variableReplacer, f.ttsService)
	default:
		return nil
	}
}
