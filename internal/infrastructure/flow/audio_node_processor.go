package flow

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// AudioNodeProcessor procesa nodos de tipo AUDIO
type AudioNodeProcessor struct {
	messagingService ports.MessagingService
	logger           ports.Logger
	variableReplacer *VariableReplacer
}

// NewAudioNodeProcessor crea un nuevo procesador de audio
func NewAudioNodeProcessor(
	messagingService ports.MessagingService,
	logger ports.Logger,
	variableReplacer *VariableReplacer,
) *AudioNodeProcessor {
	return &AudioNodeProcessor{
		messagingService: messagingService,
		logger:           logger,
		variableReplacer: variableReplacer,
	}
}

func (p *AudioNodeProcessor) Process(ctx context.Context, session *entities.FlowSession, node *entities.FlowNode) (*ProcessResult, error) {
	p.logger.Info(fmt.Sprintf("Processing AUDIO node: %s", node.ID))

	// Extraer configuración
	config := node.Config
	hasRecordedAudio, _ := config["hasRecordedAudio"].(bool)
	recordedAudio, _ := config["recordedAudio"].(string)
	waitForVoiceResponse, _ := config["waitForVoiceResponse"].(bool)
	responseVariableName, _ := config["responseVariableName"].(string)

	// CASO 1: Enviar audio al usuario
	if hasRecordedAudio && recordedAudio != "" {
		// Extraer el audio base64
		audioData := recordedAudio
		mimeType := "audio/ogg; codecs=opus" // WhatsApp usa OGG Opus
		
		if strings.HasPrefix(audioData, "data:audio/") {
			// Formato: data:audio/webm;codecs=opus;base64,UklGRiQ...
			parts := strings.Split(audioData, ",")
			if len(parts) > 1 {
				// Extraer el mime type
				headerParts := strings.Split(parts[0], ";")
				if len(headerParts) > 0 {
					mimeType = strings.TrimPrefix(headerParts[0], "data:")
				}
				audioData = parts[1]
			}
		}

		// Decodificar base64
		audioBytes, err := base64.StdEncoding.DecodeString(audioData)
		if err != nil {
			p.logger.Error(fmt.Sprintf("Error decoding audio: %v", err))
			return &ProcessResult{
				StopFlow:     true,
				ErrorMessage: fmt.Sprintf("Error processing audio: %v", err),
			}, err
		}

		p.logger.Info(fmt.Sprintf("Sending audio: %d bytes, mime: %s", len(audioBytes), mimeType))

		// Extraer número de teléfono del ConversationID (formato: phone@instance)
		phone := session.ConversationID
		if idx := strings.Index(session.ConversationID, "@"); idx != -1 {
			phone = session.ConversationID[:idx]
		}

		// Crear mensaje con audio embebido en base64
		// El MessagingService lo subirá a WhatsApp si es necesario
		message := &entities.Message{
			TenantID:       session.TenantID,
			InstanceID:     session.InstanceID,
			ConversationID: session.ConversationID,
			To:             phone, // Solo el número de teléfono
			Direction:      "out",
			MessageData: entities.MessageData{
				Type: "audio",
				Media: &entities.MediaContent{
					MimeType: mimeType,
					Data:     audioBytes,
				},
			},
		}

		err = p.messagingService.SendMessage(ctx, message)
		if err != nil {
			p.logger.Error(fmt.Sprintf("Error sending audio message: %v", err))
			return &ProcessResult{
				StopFlow:     true,
				ErrorMessage: fmt.Sprintf("Error sending message: %v", err),
			}, err
		}
		
		p.logger.Info("Audio sent successfully")

		// Si también espera respuesta de voz
		if waitForVoiceResponse {
			return &ProcessResult{
				WaitingForResponse: true,
				WaitingForVariable: responseVariableName,
				StopFlow:           false,
			}, nil
		}

		return &ProcessResult{
			WaitingForResponse: false,
			StopFlow:           false,
		}, nil
	}

	// CASO 2: Solo solicitar audio al usuario
	if waitForVoiceResponse {
		// Extraer número de teléfono del ConversationID (formato: phone@instance)
		phone := session.ConversationID
		if idx := strings.Index(session.ConversationID, "@"); idx != -1 {
			phone = session.ConversationID[:idx]
		}

		// Enviar mensaje solicitando audio
		message := &entities.Message{
			TenantID:       session.TenantID,
			InstanceID:     session.InstanceID,
			ConversationID: session.ConversationID,
			To:             phone, // Solo el número de teléfono
			Direction:      "out",
			MessageData: entities.MessageData{
				Type: "text",
				Text: &entities.TextContent{
					Body: "🎤 Por favor, envía un mensaje de voz",
				},
			},
		}

		err := p.messagingService.SendMessage(ctx, message)
		if err != nil {
			p.logger.Error(fmt.Sprintf("Error sending audio request: %v", err))
			return &ProcessResult{
				StopFlow:     true,
				ErrorMessage: fmt.Sprintf("Error sending message: %v", err),
			}, err
		}

		return &ProcessResult{
			WaitingForResponse: true,
			WaitingForVariable: responseVariableName,
			StopFlow:           false,
		}, nil
	}

	// Sin configuración válida
	return &ProcessResult{
		WaitingForResponse: false,
		StopFlow:           false,
	}, nil
}


