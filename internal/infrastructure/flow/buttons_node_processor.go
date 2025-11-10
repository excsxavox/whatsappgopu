package flow

import (
	"context"
	"fmt"
	"strings"

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
	
	// La configuración puede venir en diferentes formatos
	var header, body, footer string
	var buttonsConfig []interface{}
	var responseVariableName string

	// Formato 1: action.header, action.body, action.buttons
	if actionConfig, ok := config["action"].(map[string]interface{}); ok {
		header, _ = actionConfig["header"].(string)
		body, _ = actionConfig["body"].(string)
		footer, _ = actionConfig["footer"].(string)
		buttonsConfig, _ = actionConfig["buttons"].([]interface{})
	} else {
		// Formato 2: content/bodyText y buttons directamente
		body, _ = config["content"].(string)
		if body == "" {
			body, _ = config["bodyText"].(string)
		}
		buttonsConfig, _ = config["buttons"].([]interface{})
	}
	
	responseVariableName, _ = config["responseVariableName"].(string)

	// Reemplazar variables en el contenido
	body = p.variableReplacer.ReplaceInString(body, session.Variables)
	header = p.variableReplacer.ReplaceInString(header, session.Variables)
	footer = p.variableReplacer.ReplaceInString(footer, session.Variables)

	// Los botones ya vienen en el formato correcto de WhatsApp desde Mongo
	// Solo necesitamos copiarlos y reemplazar variables en los títulos
	buttons := []map[string]interface{}{}
	p.logger.Info(fmt.Sprintf("📋 Procesando %d botones config", len(buttonsConfig)))
	
	for i, btnConfig := range buttonsConfig {
		btnMap, ok := btnConfig.(map[string]interface{})
		if !ok {
			p.logger.Warn(fmt.Sprintf("⚠️ Botón %d no es un map[string]interface{}", i))
			continue
		}

		p.logger.Info(fmt.Sprintf("🔍 Botón %d estructura completa: %+v", i, btnMap))

		// Los botones vienen con esta estructura:
		// { "type": "reply", "reply": { "id": "...", "title": "..." } }
		// Solo necesitamos copiarlos y reemplazar variables en el título
		
		if replyData, ok := btnMap["reply"].(map[string]interface{}); ok {
			title, _ := replyData["title"].(string)
			id, _ := replyData["id"].(string)
			
			p.logger.Info(fmt.Sprintf("✅ Botón %d - id: %s, title: %s", i, id, title))
			
			// Reemplazar variables en el título
			title = p.variableReplacer.ReplaceInString(title, session.Variables)
			
			// Crear copia del botón con el título procesado
			buttons = append(buttons, map[string]interface{}{
				"type": btnMap["type"],
				"reply": map[string]interface{}{
					"id":    id,
					"title": title,
				},
			})
			p.logger.Info(fmt.Sprintf("✅ Botón %d agregado correctamente", i))
		} else {
			p.logger.Warn(fmt.Sprintf("⚠️ Botón %d no tiene estructura 'reply', se omite", i))
		}
	}
	
	p.logger.Info(fmt.Sprintf("📊 Total botones construidos: %d", len(buttons)))

	// Validar que haya botones
	if len(buttons) == 0 {
		p.logger.Error("❌ No se construyeron botones válidos")
		return &ProcessResult{
			StopFlow:     true,
			ErrorMessage: "No se pudieron construir botones válidos",
		}, fmt.Errorf("no valid buttons constructed")
	}

	// Extraer número de teléfono del ConversationID (formato: phone@instance)
	phone := session.ConversationID
	if idx := strings.Index(session.ConversationID, "@"); idx != -1 {
		phone = session.ConversationID[:idx]
	}

	// Crear mensaje interactivo con botones
	interactive := map[string]interface{}{
		"type": "button",
		"body": map[string]interface{}{
			"text": body,
		},
		"action": map[string]interface{}{
			"buttons": buttons,
		},
	}
	
	// Agregar header si existe
	if header != "" {
		interactive["header"] = map[string]interface{}{
			"type": "text",
			"text": header,
		}
	}
	
	// Agregar footer si existe
	if footer != "" {
		interactive["footer"] = map[string]interface{}{
			"text": footer,
		}
	}

	message := &entities.Message{
		TenantID:       session.TenantID,
		InstanceID:     session.InstanceID,
		ConversationID: session.ConversationID,
		To:             phone, // Solo el número de teléfono
		Direction:      "out",
		MessageData: entities.MessageData{
			Type:        "interactive",
			Interactive: interactive,
		},
	}

	// LOG: Ver estructura antes de enviar
	p.logger.Info(fmt.Sprintf("🔘 Mensaje de botones construido - To: %s, Buttons count: %d", phone, len(buttons)))
	p.logger.Info(fmt.Sprintf("   Header: %s", header))
	p.logger.Info(fmt.Sprintf("   Body: %s", body))
	p.logger.Info(fmt.Sprintf("   Footer: %s", footer))

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


