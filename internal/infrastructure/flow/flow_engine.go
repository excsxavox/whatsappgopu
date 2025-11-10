package flow

import (
	"context"
	"fmt"
	"time"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// flowEngine implementa el motor de flujos
type flowEngine struct {
	flowRepo          ports.FlowRepository
	sessionRepo       ports.FlowSessionRepository
	processorFactory  *NodeProcessorFactory
	logger            ports.Logger
	variableReplacer  *VariableReplacer
}

// NewFlowEngine crea un nuevo motor de flujos
func NewFlowEngine(
	flowRepo ports.FlowRepository,
	sessionRepo ports.FlowSessionRepository,
	messagingService ports.MessagingService,
	logger ports.Logger,
) ports.FlowEngine {
	return &flowEngine{
		flowRepo:         flowRepo,
		sessionRepo:      sessionRepo,
		processorFactory: NewNodeProcessorFactory(messagingService, logger),
		logger:           logger,
		variableReplacer: NewVariableReplacer(),
	}
}

// StartFlow inicia un nuevo flujo
func (e *flowEngine) StartFlow(ctx context.Context, conversationID string, flowID string, tenantID string, instanceID string) (*entities.FlowSession, error) {
	e.logger.Info(fmt.Sprintf("Starting flow %s for conversation %s", flowID, conversationID))

	// Buscar el flujo
	flow, err := e.flowRepo.FindByID(ctx, flowID)
	if err != nil {
		return nil, fmt.Errorf("flow not found: %w", err)
	}

	if !flow.IsActive {
		return nil, fmt.Errorf("flow %s is not active", flowID)
	}

	// Verificar si ya hay una sesión activa
	existingSession, err := e.sessionRepo.FindActiveByConversation(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("error checking existing session: %w", err)
	}

	if existingSession != nil {
		e.logger.Warn(fmt.Sprintf("Active session already exists for conversation %s", conversationID))
		return existingSession, nil
	}

	// Crear nueva sesión
	session := entities.NewFlowSession(conversationID, flowID, flow.GetEntryNodeID(), tenantID, instanceID)

	// Guardar sesión
	err = e.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("error saving session: %w", err)
	}

	// Procesar el nodo de entrada
	entryNode := flow.GetNodeByID(flow.GetEntryNodeID())
	if entryNode == nil {
		return nil, fmt.Errorf("entry node %s not found in flow", flow.GetEntryNodeID())
	}

	err = e.ProcessNode(ctx, session, entryNode)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Error processing entry node: %v", err))
		session.MarkError()
		e.sessionRepo.Update(ctx, session)
		return nil, err
	}

	return session, nil
}

// ProcessMessage procesa un mensaje en el contexto de un flujo
func (e *flowEngine) ProcessMessage(ctx context.Context, session *entities.FlowSession, message *entities.Message) error {
	e.logger.Info(fmt.Sprintf("Processing message in flow session %s", session.ID))

	// Actualizar última actividad
	session.UpdateActivity()

	// Buscar el flujo
	_, err := e.flowRepo.FindByID(ctx, session.FlowID)
	if err != nil {
		return fmt.Errorf("flow not found: %w", err)
	}

	// Si está esperando respuesta, capturar el valor
	if session.WaitingForResponse {
		// Extraer valor según tipo de mensaje
		var value interface{}

		switch message.MessageData.Type {
		case "text":
			if message.MessageData.Text != nil {
				value = message.MessageData.Text.Body
			}
		case "image":
			if message.MessageData.Media != nil && message.MessageData.Media.Storage != nil {
				// Preferir la URL pública, o la key si no hay URL
				if message.MessageData.Media.Storage.PublicURL != "" {
					value = message.MessageData.Media.Storage.PublicURL
				} else {
					value = message.MessageData.Media.Storage.Key
				}
			}
		case "audio":
			if message.MessageData.Media != nil && message.MessageData.Media.Storage != nil {
				if message.MessageData.Media.Storage.PublicURL != "" {
					value = message.MessageData.Media.Storage.PublicURL
				} else {
					value = message.MessageData.Media.Storage.Key
				}
			}
		case "interactive":
			// Botón presionado
			if message.MessageData.Interactive != nil {
				if buttonReply, ok := message.MessageData.Interactive["button_reply"].(map[string]interface{}); ok {
					if id, ok := buttonReply["id"].(string); ok && id != "" {
						value = id
					}
				} else if listReply, ok := message.MessageData.Interactive["list_reply"].(map[string]interface{}); ok {
					if id, ok := listReply["id"].(string); ok && id != "" {
						value = id
					}
				}
			}
		default:
			e.logger.Warn(fmt.Sprintf("Unsupported message type: %s", message.MessageData.Type))
			return fmt.Errorf("unsupported message type: %s", message.MessageData.Type)
		}

		// Guardar en variable
		if session.WaitingForVariable != "" {
			session.SetVariable(session.WaitingForVariable, value)
			e.logger.Info(fmt.Sprintf("Captured variable %s = %v", session.WaitingForVariable, value))
		}

		// Limpiar estado de espera
		session.ClearWaitingForResponse()

		// Actualizar sesión
		err = e.sessionRepo.Update(ctx, session)
		if err != nil {
			return fmt.Errorf("error updating session: %w", err)
		}

		// Avanzar al siguiente nodo
		err = e.MoveToNextNode(ctx, session, "default")
		if err != nil {
			return fmt.Errorf("error moving to next node: %w", err)
		}

		return nil
	}

	// Si no está esperando respuesta, procesar como nuevo mensaje
	// (esto podría ser un comando o reiniciar el flujo)
	e.logger.Warn("Message received but not waiting for response")
	return nil
}

// ProcessNode procesa un nodo específico
func (e *flowEngine) ProcessNode(ctx context.Context, session *entities.FlowSession, node *entities.FlowNode) error {
	e.logger.Info(fmt.Sprintf("Processing node %s (type: %s)", node.ID, node.Type))

	// Obtener procesador para este tipo de nodo
	processor := e.processorFactory.GetProcessor(node.Type)
	if processor == nil {
		return fmt.Errorf("no processor found for node type: %s", node.Type)
	}

	// Procesar el nodo
	result, err := processor.Process(ctx, session, node)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Error processing node: %v", err))
		return err
	}

	// Actualizar sesión según resultado
	if result.StopFlow {
		e.logger.Info("Flow stopped by node processor")
		session.MarkError()
		return e.sessionRepo.Update(ctx, session)
	}

	if result.WaitingForResponse {
		e.logger.Info(fmt.Sprintf("Waiting for response, variable: %s", result.WaitingForVariable))
		session.SetWaitingForResponse(result.WaitingForVariable)
		return e.sessionRepo.Update(ctx, session)
	}

	// Si es un nodo CONDITION, result.NextNodeID contiene "yes" o "no"
	condition := "default"
	if node.Type == "CONDITION" {
		condition = result.NextNodeID
	}

	// Si no espera respuesta, avanzar al siguiente nodo
	err = e.MoveToNextNode(ctx, session, condition)
	if err != nil {
		return fmt.Errorf("error moving to next node: %w", err)
	}

	return nil
}

// MoveToNextNode avanza al siguiente nodo según los edges
func (e *flowEngine) MoveToNextNode(ctx context.Context, session *entities.FlowSession, condition string) error {
	e.logger.Info(fmt.Sprintf("Moving to next node from %s with condition: %s", session.CurrentNodeID, condition))

	// Buscar el flujo
	flow, err := e.flowRepo.FindByID(ctx, session.FlowID)
	if err != nil {
		return fmt.Errorf("flow not found: %w", err)
	}

	// Buscar edges salientes del nodo actual
	outgoingEdges := flow.GetOutgoingEdges(session.CurrentNodeID)

	if len(outgoingEdges) == 0 {
		e.logger.Info("No outgoing edges, completing flow")
		session.Complete()
		return e.sessionRepo.Update(ctx, session)
	}

	// Seleccionar el edge apropiado
	var selectedEdge *entities.FlowEdge

	if condition != "default" {
		// Buscar edge con la condición específica
		for i := range outgoingEdges {
			if outgoingEdges[i].Condition == condition || 
			   (condition == "yes" && (outgoingEdges[i].Condition == "si" || outgoingEdges[i].Condition == "yes")) ||
			   (condition == "no" && outgoingEdges[i].Condition == "no") {
				selectedEdge = &outgoingEdges[i]
				break
			}
		}
	} else {
		// Tomar el primer edge disponible
		selectedEdge = &outgoingEdges[0]
	}

	if selectedEdge == nil {
		e.logger.Error(fmt.Sprintf("No edge found for condition: %s", condition))
		session.Complete()
		return e.sessionRepo.Update(ctx, session)
	}

	// Aplicar delay si existe
	if selectedEdge.DelayMs > 0 {
		time.Sleep(time.Duration(selectedEdge.DelayMs) * time.Millisecond)
	}

	// Mover al siguiente nodo
	nextNodeID := selectedEdge.To
	nextNode := flow.GetNodeByID(nextNodeID)
	if nextNode == nil {
		return fmt.Errorf("next node %s not found", nextNodeID)
	}

	session.MoveToNode(nextNodeID)

	// Actualizar sesión
	err = e.sessionRepo.Update(ctx, session)
	if err != nil {
		return fmt.Errorf("error updating session: %w", err)
	}

	// Procesar el siguiente nodo
	return e.ProcessNode(ctx, session, nextNode)
}


