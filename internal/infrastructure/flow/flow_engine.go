package flow

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// flowEngine implementa el motor de flujos
type flowEngine struct {
	flowRepo             ports.FlowRepository
	sessionRepo          ports.FlowSessionRepository
	processorFactory     *NodeProcessorFactory
	logger               ports.Logger
	variableReplacer     *VariableReplacer
	transcriptionService TranscriptionService
	messagingService     ports.MessagingService
}

// TranscriptionService interface for audio transcription
type TranscriptionService interface {
	TranscribeAudio(audioData []byte, filename string) (string, error)
}

// NewFlowEngine crea un nuevo motor de flujos
func NewFlowEngine(
	flowRepo ports.FlowRepository,
	sessionRepo ports.FlowSessionRepository,
	messagingService ports.MessagingService,
	logger ports.Logger,
	transcriptionService TranscriptionService,
	aiValidationService AIValidationService,
	aiContextService AIContextService,
	contextoRepo ports.ContextoRepository,
	ttsService TTSService,
) ports.FlowEngine {
	return &flowEngine{
		flowRepo:             flowRepo,
		sessionRepo:          sessionRepo,
		processorFactory:     NewNodeProcessorFactory(messagingService, logger, aiValidationService, aiContextService, contextoRepo, flowRepo, ttsService),
		logger:               logger,
		variableReplacer:     NewVariableReplacer(),
		transcriptionService: transcriptionService,
		messagingService:     messagingService,
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
			if message.MessageData.Media != nil {
				// Intentar transcribir el audio
				transcribedText, err := e.transcribeAudioMessage(message)
				if err != nil {
					e.logger.Error(fmt.Sprintf("Error transcribing audio: %v", err))
					// Si falla la transcripción, guardar la URL como fallback
					if message.MessageData.Media.Storage != nil {
						if message.MessageData.Media.Storage.PublicURL != "" {
							value = message.MessageData.Media.Storage.PublicURL
						} else {
							value = message.MessageData.Media.Storage.Key
						}
					}
				} else {
					// Guardar el texto transcrito
					value = transcribedText
					e.logger.Info(fmt.Sprintf("🎤 Audio transcrito: %s", transcribedText))
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
			// También guardar el tipo de mensaje para validación posterior
			session.SetVariable(session.WaitingForVariable+"_message_type", message.MessageData.Type)
			// Limpiar el flag de validación cuando se captura una nueva respuesta
			// Esto permite que la nueva respuesta sea validada
			session.SetVariable(session.WaitingForVariable+"_validated", false)
			e.logger.Info(fmt.Sprintf("Captured variable %s = %v (type: %s)", session.WaitingForVariable, value, message.MessageData.Type))
		}

		// Limpiar estado de espera
		session.ClearWaitingForResponse()

		// Actualizar sesión
		err = e.sessionRepo.Update(ctx, session)
		if err != nil {
			return fmt.Errorf("error updating session: %w", err)
		}

		// Buscar el flujo y el nodo actual
		flow, err := e.flowRepo.FindByID(ctx, session.FlowID)
		if err != nil {
			return fmt.Errorf("error finding flow: %w", err)
		}

		var currentNode *entities.FlowNode
		for i := range flow.FlowData.Nodes {
			if flow.FlowData.Nodes[i].ID == session.CurrentNodeID {
				currentNode = &flow.FlowData.Nodes[i]
				break
			}
		}

		// Si el nodo actual es RESPONSE, procesarlo para validar la variable
		// Si es TEXT, BUTTONS, AUDIO, etc., avanzar al siguiente nodo (que debería ser RESPONSE)
		if currentNode != nil && currentNode.Type == "RESPONSE" {
			e.logger.Info(fmt.Sprintf("Processing RESPONSE node %s after capturing variable", currentNode.ID))

			// Verificar si este nodo RESPONSE ya validó exitosamente (evitar bucles)
			config := currentNode.Config
			variableName, _ := config["variableName"].(string)
			if variableName != "" {
				validatedKey := variableName + "_validated"
				if validated, exists := session.GetVariable(validatedKey); exists {
					if validatedBool, ok := validated.(bool); ok && validatedBool {
						e.logger.Info(fmt.Sprintf("⚠️ Variable %s ya fue validada exitosamente en este nodo, evitando bucle. Avanzando al siguiente nodo.", variableName))
						// Limpiar el flag de validación para permitir nueva validación si el usuario envía nueva respuesta
						session.SetVariable(validatedKey, false)
						// Avanzar al siguiente nodo sin reprocesar
						err = e.MoveToNextNode(ctx, session, "default")
						if err != nil {
							return fmt.Errorf("error moving to next node: %w", err)
						}
						return nil
					}
				}
			}

			err = e.ProcessNode(ctx, session, currentNode)
			if err != nil {
				return fmt.Errorf("error processing RESPONSE node: %w", err)
			}
		} else {
			// Para nodos TEXT, BUTTONS, AUDIO, etc., avanzar al siguiente nodo
			nodeType := "unknown"
			if currentNode != nil {
				nodeType = currentNode.Type
			}
			e.logger.Info(fmt.Sprintf("Current node is %s (type: %s), moving to next node for validation",
				session.CurrentNodeID, nodeType))
			err = e.MoveToNextNode(ctx, session, "default")
			if err != nil {
				return fmt.Errorf("error moving to next node: %w", err)
			}
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
			edge := &outgoingEdges[i]
			edgeCondition := edge.GetCondition()

			// Verificar por Condition o por SourceHandle
			matchCondition := edgeCondition == condition ||
				(condition == "yes" && (edgeCondition == "si" || edgeCondition == "yes")) ||
				(condition == "no" && edgeCondition == "no")

			// Si no hay Condition, usar SourceHandle (default = yes, output = no)
			matchSourceHandle := (condition == "yes" && edge.SourceHandle == "default") ||
				(condition == "no" && edge.SourceHandle == "output")

			if matchCondition || matchSourceHandle {
				e.logger.Info(fmt.Sprintf("📍 Selected edge: %s -> %s (condition: %s, sourceHandle: %s)",
					edge.From, edge.To, edgeCondition, edge.SourceHandle))
				selectedEdge = edge
				break
			}
		}
	} else {
		// Tomar el primer edge disponible (o el que tiene sourceHandle "default")
		// PERO evitar edges que van a nodos de error si la validación pasó
		// Priorizar edges que NO van a nodos de error
		var errorEdges []*entities.FlowEdge
		var normalEdges []*entities.FlowEdge

		for i := range outgoingEdges {
			edge := &outgoingEdges[i]

			// Verificar si el nodo destino es un nodo de error
			nextNode := flow.GetNodeByID(edge.To)
			isErrorNode := false
			if nextNode != nil {
				config := nextNode.Config
				isErrorNode, _ = config["isErrorNode"].(bool)
				errorForResponseNode, _ := config["errorForResponseNode"].(string)

				// Si es un nodo de error y está asociado al nodo actual
				if isErrorNode && errorForResponseNode == session.CurrentNodeID {
					errorEdges = append(errorEdges, edge)
					continue
				}
			}

			// Edge normal (no es de error)
			normalEdges = append(normalEdges, edge)
		}

		// Si hay edges normales, usar uno de esos (evitar nodos de error cuando la validación pasa)
		if len(normalEdges) > 0 {
			for _, edge := range normalEdges {
				if edge.SourceHandle == "default" || selectedEdge == nil {
					selectedEdge = edge
					if edge.SourceHandle == "default" {
						break
					}
				}
			}
		} else if len(errorEdges) > 0 {
			// Solo usar edge de error si no hay edges normales disponibles
			e.logger.Warn(fmt.Sprintf("⚠️ Solo hay edges a nodos de error disponibles, usando el primero"))
			selectedEdge = errorEdges[0]
		}
	}

	if selectedEdge == nil {
		e.logger.Error(fmt.Sprintf("❌ No edge found for condition: %s. Available edges:", condition))
		for _, edge := range outgoingEdges {
			e.logger.Error(fmt.Sprintf("  - Edge: %s -> %s (condition: %s, sourceHandle: %s)",
				edge.From, edge.To, edge.GetCondition(), edge.SourceHandle))
		}
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

// transcribeAudioMessage descarga y transcribe un mensaje de audio
func (e *flowEngine) transcribeAudioMessage(message *entities.Message) (string, error) {
	if message.MessageData.Media == nil {
		return "", fmt.Errorf("no media data in message")
	}

	var audioData []byte
	var err error

	// Caso 1: Si ya tenemos el contenido del audio en memoria
	if message.MessageData.Media.Data != nil && len(message.MessageData.Media.Data) > 0 {
		e.logger.Info("📦 Usando audio desde Data en memoria")
		audioData = message.MessageData.Media.Data
	} else if message.MessageData.Media.Storage != nil {
		// Caso 2: Descargar desde WhatsApp usando media_id
		if message.MessageData.Media.Storage.Provider == "whatsapp" && message.MessageData.Media.Storage.Key != "" {
			e.logger.Info(fmt.Sprintf("📥 Descargando audio desde WhatsApp, media_id: %s", message.MessageData.Media.Storage.Key))
			audioData, err = e.downloadAudioFromWhatsApp(message.MessageData.Media.Storage.Key)
			if err != nil {
				return "", fmt.Errorf("error downloading audio from WhatsApp: %w", err)
			}
		} else if message.MessageData.Media.Storage.PublicURL != "" {
			// Caso 3: Descargar desde URL pública
			e.logger.Info(fmt.Sprintf("📥 Descargando audio desde URL: %s", message.MessageData.Media.Storage.PublicURL))
			audioData, err = e.downloadAudioFromURL(message.MessageData.Media.Storage.PublicURL)
			if err != nil {
				return "", fmt.Errorf("error downloading audio from URL: %w", err)
			}
		} else {
			return "", fmt.Errorf("no valid audio source in Storage (provider: %s, key: %s)",
				message.MessageData.Media.Storage.Provider, message.MessageData.Media.Storage.Key)
		}
	} else {
		return "", fmt.Errorf("no audio source available (no Data or Storage)")
	}

	// Transcribir el audio
	if e.transcriptionService == nil {
		return "", fmt.Errorf("transcription service not configured")
	}

	filename := "audio.ogg"
	if message.MessageData.Media.MimeType != "" {
		// Determinar extensión basada en mime type
		if message.MessageData.Media.MimeType == "audio/ogg" {
			filename = "audio.ogg"
		} else if message.MessageData.Media.MimeType == "audio/mpeg" {
			filename = "audio.mp3"
		} else if message.MessageData.Media.MimeType == "audio/mp4" || message.MessageData.Media.MimeType == "audio/aac" {
			filename = "audio.m4a"
		}
	}

	return e.transcriptionService.TranscribeAudio(audioData, filename)
}

// downloadAudioFromWhatsApp descarga audio desde WhatsApp usando media_id
func (e *flowEngine) downloadAudioFromWhatsApp(mediaID string) ([]byte, error) {
	data, mimeType, err := e.messagingService.DownloadMedia(mediaID)
	if err != nil {
		return nil, fmt.Errorf("error downloading from WhatsApp: %w", err)
	}

	e.logger.Info(fmt.Sprintf("✅ Audio descargado de WhatsApp: %d bytes, mime: %s", len(data), mimeType))
	return data, nil
}

// downloadAudioFromURL descarga audio desde una URL pública
func (e *flowEngine) downloadAudioFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading audio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error downloading audio, status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
