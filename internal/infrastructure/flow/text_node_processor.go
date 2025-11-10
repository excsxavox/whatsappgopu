package flow

import (
	"context"
	"fmt"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// TextNodeProcessor procesa nodos de tipo TEXT
type TextNodeProcessor struct {
	messagingService ports.MessagingService
	logger           ports.Logger
	variableReplacer *VariableReplacer
}

// NewTextNodeProcessor crea un nuevo procesador de texto
func NewTextNodeProcessor(
	messagingService ports.MessagingService,
	logger ports.Logger,
	variableReplacer *VariableReplacer,
) *TextNodeProcessor {
	return &TextNodeProcessor{
		messagingService: messagingService,
		logger:           logger,
		variableReplacer: variableReplacer,
	}
}

func (p *TextNodeProcessor) Process(ctx context.Context, session *entities.FlowSession, node *entities.FlowNode) (*ProcessResult, error) {
	p.logger.Info(fmt.Sprintf("Processing TEXT node: %s", node.ID))

	// Extraer configuración
	config := node.Config
	content, _ := config["content"].(string)
	if content == "" {
		content, _ = config["bodyText"].(string)
	}

	waitForResponse, _ := config["waitForResponse"].(bool)
	responseVariableName, _ := config["responseVariableName"].(string)

	// Reemplazar variables en el contenido
	content = p.variableReplacer.ReplaceInString(content, session.Variables)

	// Crear mensaje
	message := &entities.Message{
		TenantID:       session.TenantID,
		InstanceID:     session.InstanceID,
		ConversationID: session.ConversationID,
		From:           session.ConversationID, // Número del usuario
		Type:           "text",
		Text: entities.MessageText{
			Body: content,
		},
	}

	// Enviar mensaje
	err := p.messagingService.SendMessage(ctx, message)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Error sending TEXT message: %v", err))
		return &ProcessResult{
			StopFlow:     true,
			ErrorMessage: fmt.Sprintf("Error sending message: %v", err),
		}, err
	}

	// Determinar si debe esperar respuesta
	result := &ProcessResult{
		WaitingForResponse: waitForResponse,
		WaitingForVariable: responseVariableName,
		StopFlow:           false,
	}

	return result, nil
}


