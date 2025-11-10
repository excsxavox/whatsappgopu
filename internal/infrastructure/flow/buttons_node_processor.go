package flow

import (
	"context"
	"fmt"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// ButtonsNodeProcessor procesa nodos de tipo BUTTONS
type ButtonsNodeProcessor struct {
	messagingService ports.MessagingService
	logger           ports.Logger
	variableReplacer *VariableReplacer
}

// NewButtonsNodeProcessor crea un nuevo procesador de botones
func NewButtonsNodeProcessor(
	messagingService ports.MessagingService,
	logger ports.Logger,
	variableReplacer *VariableReplacer,
) *ButtonsNodeProcessor {
	return &ButtonsNodeProcessor{
		messagingService: messagingService,
		logger:           logger,
		variableReplacer: variableReplacer,
	}
}

func (p *ButtonsNodeProcessor) Process(ctx context.Context, session *entities.FlowSession, node *entities.FlowNode) (*ProcessResult, error) {
	p.logger.Info(fmt.Sprintf("Processing BUTTONS node: %s", node.ID))

	// Extraer configuración
	config := node.Config
	content, _ := config["content"].(string)
	buttonsConfig, _ := config["buttons"].([]interface{})
	responseVariableName, _ := config["responseVariableName"].(string)

	// Reemplazar variables en el contenido
	content = p.variableReplacer.ReplaceInString(content, session.Variables)

	// Construir botones para el payload (no usar entidades, usar map)
	buttons := []map[string]interface{}{}
	for _, btnConfig := range buttonsConfig {
		btnMap, ok := btnConfig.(map[string]interface{})
		if !ok {
			continue
		}

		btnID, _ := btnMap["id"].(string)
		btnTitle, _ := btnMap["title"].(string)

		// Reemplazar variables en título
		btnTitle = p.variableReplacer.ReplaceInString(btnTitle, session.Variables)

		buttons = append(buttons, map[string]interface{}{
			"type": "reply",
			"reply": map[string]interface{}{
				"id":    btnID,
				"title": btnTitle,
			},
		})
	}

	// Crear mensaje interactivo
	message := &entities.Message{
		TenantID:       session.TenantID,
		InstanceID:     session.InstanceID,
		ConversationID: session.ConversationID,
		To:             session.ConversationID,
		Direction:      "out",
		MessageData: entities.MessageData{
			Type: "interactive",
			Interactive: &entities.InteractiveContent{
				Type: "button",
			},
		},
	}

	// Agregar los botones como metadata adicional (el adapter los manejará)
	// Por ahora, enviamos un mensaje de texto con los botones como fallback
	message.MessageData.Type = "text"
	message.MessageData.Text = &entities.TextContent{
		Body: content,
	}

	// Enviar mensaje
	err := p.messagingService.SendMessage(ctx, message)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Error sending BUTTONS message: %v", err))
		return &ProcessResult{
			StopFlow:     true,
			ErrorMessage: fmt.Sprintf("Error sending message: %v", err),
		}, err
	}

	// Los botones siempre esperan respuesta
	result := &ProcessResult{
		WaitingForResponse: true,
		WaitingForVariable: responseVariableName,
		StopFlow:           false,
	}

	return result, nil
}


