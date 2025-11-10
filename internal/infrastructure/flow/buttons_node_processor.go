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

	// Construir botones
	buttons := []entities.InteractiveButton{}
	for _, btnConfig := range buttonsConfig {
		btnMap, ok := btnConfig.(map[string]interface{})
		if !ok {
			continue
		}

		btnID, _ := btnMap["id"].(string)
		btnType, _ := btnMap["type"].(string)
		btnTitle, _ := btnMap["title"].(string)

		// Reemplazar variables en título
		btnTitle = p.variableReplacer.ReplaceInString(btnTitle, session.Variables)

		buttons = append(buttons, entities.InteractiveButton{
			Type: btnType,
			Reply: entities.InteractiveReply{
				ID:    btnID,
				Title: btnTitle,
			},
		})
	}

	// Crear mensaje interactivo
	message := &entities.Message{
		TenantID:       session.TenantID,
		InstanceID:     session.InstanceID,
		ConversationID: session.ConversationID,
		From:           session.ConversationID,
		Type:           "interactive",
		Interactive: entities.MessageInteractive{
			Type: "button",
			Body: entities.InteractiveBody{
				Text: content,
			},
			Action: entities.InteractiveAction{
				Buttons: buttons,
			},
		},
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


