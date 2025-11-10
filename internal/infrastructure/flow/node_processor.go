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

// NodeProcessorFactory crea procesadores según el tipo de nodo
type NodeProcessorFactory struct {
	messagingService ports.MessagingService
	logger           ports.Logger
	variableReplacer *VariableReplacer
}

// NewNodeProcessorFactory crea una nueva factory
func NewNodeProcessorFactory(
	messagingService ports.MessagingService,
	logger ports.Logger,
) *NodeProcessorFactory {
	return &NodeProcessorFactory{
		messagingService: messagingService,
		logger:           logger,
		variableReplacer: NewVariableReplacer(),
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
		return NewResponseNodeProcessor(f.logger, f.variableReplacer)
	case "AUDIO":
		return NewAudioNodeProcessor(f.messagingService, f.logger, f.variableReplacer)
	default:
		return nil
	}
}


