package flow

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// AIValidationService interface for AI-based validation
type AIValidationService interface {
	ValidateWithAI(userResponse, validationPrompt string) (bool, string, error)
}

// ResponseNodeProcessor procesa nodos de tipo RESPONSE
type ResponseNodeProcessor struct {
	logger              ports.Logger
	variableReplacer    *VariableReplacer
	aiValidationService AIValidationService
	aiContextService    AIContextService
	contextoRepo        ports.ContextoRepository
	flowRepo            ports.FlowRepository
	messagingService    ports.MessagingService
}

// NewResponseNodeProcessor crea un nuevo procesador de respuestas
func NewResponseNodeProcessor(
	logger ports.Logger,
	variableReplacer *VariableReplacer,
	aiValidationService AIValidationService,
	aiContextService AIContextService,
	contextoRepo ports.ContextoRepository,
	flowRepo ports.FlowRepository,
	messagingService ports.MessagingService,
) *ResponseNodeProcessor {
	return &ResponseNodeProcessor{
		logger:              logger,
		variableReplacer:    variableReplacer,
		aiValidationService: aiValidationService,
		aiContextService:    aiContextService,
		contextoRepo:        contextoRepo,
		flowRepo:            flowRepo,
		messagingService:    messagingService,
	}
}

func (p *ResponseNodeProcessor) Process(ctx context.Context, session *entities.FlowSession, node *entities.FlowNode) (*ProcessResult, error) {
	p.logger.Info(fmt.Sprintf("Processing RESPONSE node: %s", node.ID))

	// Extraer configuración
	config := node.Config
	variableName, _ := config["variableName"].(string)
	validationConfig, _ := config["validation"].(map[string]interface{})

	// Verificar si esta variable ya fue validada exitosamente (evitar bucles)
	validatedKey := variableName + "_validated"
	if validated, exists := session.GetVariable(validatedKey); exists {
		if validatedBool, ok := validated.(bool); ok && validatedBool {
			p.logger.Info(fmt.Sprintf("⚠️ Variable %s ya fue validada exitosamente, evitando bucle. Avanzando al siguiente nodo.", variableName))
			// Si ya fue validada, simplemente avanzar sin reprocesar
			return &ProcessResult{
				WaitingForResponse: false,
				StopFlow:           false,
			}, nil
		}
	}

	// Obtener valor de la variable
	value, exists := session.GetVariable(variableName)
	if !exists {
		p.logger.Warn(fmt.Sprintf("Variable %s not found in session", variableName))
		return &ProcessResult{
			StopFlow:     false,
			ErrorMessage: fmt.Sprintf("Variable %s not found", variableName),
		}, nil
	}

	valueStr := fmt.Sprintf("%v", value)

	// Log del valor capturado para debugging
	p.logger.Info(fmt.Sprintf("🔍 Valor capturado para variable %s: '%s' (longitud: %d)", variableName, valueStr, len(valueStr)))

	// 📋 VALIDACIÓN DE TIPO DE MENSAJE (texto vs audio vs otros)
	responseType, _ := config["responseType"].(string)
	if responseType != "" {
		// Obtener el tipo real del mensaje que se capturó
		actualMessageType, typeExists := session.GetVariable(variableName + "_message_type")
		actualTypeStr := ""
		if typeExists {
			actualTypeStr = fmt.Sprintf("%v", actualMessageType)
		}

		p.logger.Info(fmt.Sprintf("🔍 Validación de tipo: esperado=%s, recibido=%s", responseType, actualTypeStr))

		// Verificar si el tipo coincide
		typeMatch := false
		switch responseType {
		case "text":
			typeMatch = (actualTypeStr == "text" || actualTypeStr == "interactive")
		case "audio":
			typeMatch = (actualTypeStr == "audio")
		case "image":
			typeMatch = (actualTypeStr == "image")
		case "video":
			typeMatch = (actualTypeStr == "video")
		case "document":
			typeMatch = (actualTypeStr == "document")
		default:
			typeMatch = true // Si no se especifica tipo, aceptar cualquiera
		}

		if !typeMatch {
			p.logger.Warn(fmt.Sprintf("❌ Tipo de mensaje incorrecto: esperado %s, recibido %s", responseType, actualTypeStr))

			// Determinar mensaje de error según el tipo esperado
			var typeErrorMessage string
			switch responseType {
			case "text":
				typeErrorMessage = "Por favor envía un mensaje de texto, no un audio o archivo."
			case "audio":
				typeErrorMessage = "Por favor envía un mensaje de voz."
			case "image":
				typeErrorMessage = "Por favor envía una imagen."
			case "video":
				typeErrorMessage = "Por favor envía un video."
			case "document":
				typeErrorMessage = "Por favor envía un documento."
			default:
				typeErrorMessage = "Tipo de mensaje no válido."
			}

			// Enviar mensaje de error al usuario
			if p.messagingService != nil {
				phone := session.ConversationID
				if idx := len(phone) - 1; idx > 0 {
					for i, ch := range phone {
						if ch == '@' {
							phone = phone[:i]
							break
						}
					}
				}

				errorMsg := &entities.Message{
					TenantID:       session.TenantID,
					InstanceID:     session.InstanceID,
					ConversationID: session.ConversationID,
					To:             phone,
					Direction:      "out",
					MessageData: entities.MessageData{
						Type: "text",
						Text: &entities.TextContent{
							Body: typeErrorMessage,
						},
					},
				}

				if err := p.messagingService.SendMessage(ctx, errorMsg); err != nil {
					p.logger.Error(fmt.Sprintf("Error enviando mensaje de error de tipo: %v", err))
				} else {
					p.logger.Info("📤 Mensaje de error de tipo enviado al usuario")
				}
			}

			// Limpiar el flag de validación para permitir nueva validación
			session.SetVariable(variableName+"_validated", false)
			// Marcar que sigue esperando respuesta (no avanzar)
			session.SetWaitingForResponse(variableName)

			return &ProcessResult{
				WaitingForResponse: true,
				WaitingForVariable: variableName,
				StopFlow:           false,
				ErrorMessage:       typeErrorMessage,
			}, nil
		}

		p.logger.Info(fmt.Sprintf("✅ Validación de tipo exitosa: %s", actualTypeStr))
	}

	// 🤖 VALIDACIÓN CON IA (si está habilitada)
	validateWithAI, _ := config["validateWithAI"].(bool)
	aiValidationPrompt, _ := config["aiValidationPrompt"].(string)
	errorMessage, _ := config["errorMessage"].(string)

	// Log de configuración de validación IA
	p.logger.Info(fmt.Sprintf("🔍 [Validación IA] Configuración - validateWithAI: %v, aiValidationPrompt: '%s' (longitud: %d), errorMessage: '%s' (longitud: %d), aiValidationService disponible: %v",
		validateWithAI, aiValidationPrompt, len(aiValidationPrompt), errorMessage, len(errorMessage), p.aiValidationService != nil))

	if validateWithAI {
		if p.aiValidationService == nil {
			p.logger.Warn("⚠️ [Validación IA] validateWithAI está habilitado pero aiValidationService es nil")
		} else if aiValidationPrompt == "" {
			p.logger.Warn("⚠️ [Validación IA] validateWithAI está habilitado pero aiValidationPrompt está vacío")
		} else {
			// Validación con IA está configurada correctamente
			p.logger.Info(fmt.Sprintf("✅ [Validación IA] Configuración correcta, procediendo con validación"))
			// Verificar si el valor está vacío antes de validar
			if valueStr == "" || len(strings.TrimSpace(valueStr)) == 0 {
				p.logger.Warn(fmt.Sprintf("⚠️ Valor vacío para variable %s, no se puede validar con IA", variableName))
				// Si el valor está vacío, tratar como inválido
				// PERO si errorMessage está vacío, usar contexto para generar respuesta
				finalErrorMessage := errorMessage
				if finalErrorMessage == "" {
					// Intentar obtener contexto y usar IA para generar respuesta
					contexto, err := p.getContextoForSession(ctx, session)
					if err != nil {
						p.logger.Warn(fmt.Sprintf("⚠️ No se pudo obtener contexto: %v, usando mensaje por defecto", err))
						finalErrorMessage = "No se pudo procesar tu mensaje. Por favor intenta de nuevo."
					} else if contexto != nil && p.aiContextService != nil {
						// Hay contexto y servicio de IA disponible
						p.logger.Info("🤖 Usando IA con contexto para generar respuesta (valor vacío)")
						// Usar un mensaje descriptivo para la IA sobre el problema
						userQuestion := fmt.Sprintf("El usuario envió un mensaje de audio pero no se pudo transcribir correctamente. El texto está vacío. ¿Cómo debo responderle de manera amigable para pedirle que intente de nuevo?")
						aiResponse, err := p.aiContextService.RespondWithContext(userQuestion, contexto)
						if err == nil && aiResponse != "" {
							finalErrorMessage = aiResponse
							p.logger.Info(fmt.Sprintf("✅ Respuesta generada por IA con contexto: %s", aiResponse))
						} else {
							p.logger.Warn(fmt.Sprintf("⚠️ Error generando respuesta con IA: %v, usando mensaje por defecto", err))
							finalErrorMessage = "No se pudo procesar tu mensaje. Por favor intenta de nuevo."
						}
					} else {
						// No hay contexto o servicio de IA
						if contexto == nil {
							p.logger.Info("ℹ️ El flow no tiene idcontext, no se usará IA con contexto")
						}
						if p.aiContextService == nil {
							p.logger.Warn("⚠️ Servicio de IA con contexto no disponible")
						}
						finalErrorMessage = "No se pudo procesar tu mensaje. Por favor intenta de nuevo."
					}
				}

				// Enviar mensaje de error al usuario
				if p.messagingService != nil {
					phone := session.ConversationID
					if idx := len(phone) - 1; idx > 0 {
						for i, ch := range phone {
							if ch == '@' {
								phone = phone[:i]
								break
							}
						}
					}

					errorMsg := &entities.Message{
						TenantID:       session.TenantID,
						InstanceID:     session.InstanceID,
						ConversationID: session.ConversationID,
						To:             phone,
						Direction:      "out",
						MessageData: entities.MessageData{
							Type: "text",
							Text: &entities.TextContent{
								Body: finalErrorMessage,
							},
						},
					}

					if err := p.messagingService.SendMessage(ctx, errorMsg); err != nil {
						p.logger.Error(fmt.Sprintf("Error enviando mensaje de error: %v", err))
					} else {
						p.logger.Info("📤 Mensaje de error enviado al usuario (valor vacío)")
					}
				}

				session.SetVariable(variableName+"_validated", false)
				session.SetWaitingForResponse(variableName)

				return &ProcessResult{
					WaitingForResponse: true,
					WaitingForVariable: variableName,
					StopFlow:           false,
					ErrorMessage:       finalErrorMessage,
				}, nil
			}

			p.logger.Info(fmt.Sprintf("🤖 [Validación IA] Iniciando validación para variable %s con valor: '%s' (longitud: %d)", variableName, valueStr, len(valueStr)))
			p.logger.Info(fmt.Sprintf("📝 [Validación IA] Prompt de validación: '%s'", aiValidationPrompt))

			isValid, reason, err := p.aiValidationService.ValidateWithAI(valueStr, aiValidationPrompt)
			p.logger.Info(fmt.Sprintf("📊 [Validación IA] Resultado - isValid: %v, reason: '%s', error: %v", isValid, reason, err))
			if err != nil {
				p.logger.Error(fmt.Sprintf("Error en validación IA: %v", err))
				// Continuar con validaciones normales si IA falla
			} else if !isValid {
				p.logger.Warn(fmt.Sprintf("❌ Validación IA falló: %s", reason))

				// Si errorMessage está vacío, usar IA con contexto para generar respuesta
				finalErrorMessage := errorMessage
				if finalErrorMessage == "" {
					// Intentar obtener contexto y usar IA
					contexto, err := p.getContextoForSession(ctx, session)
					if err != nil {
						// Error al buscar contexto (ej: flow no encontrado)
						p.logger.Warn(fmt.Sprintf("⚠️ No se pudo obtener contexto: %v, usando mensaje por defecto", err))
						finalErrorMessage = "La respuesta no es válida. Por favor intenta de nuevo."
					} else if contexto != nil && p.aiContextService != nil {
						// Hay contexto y servicio de IA disponible
						p.logger.Info("🤖 Usando IA con contexto para generar respuesta de error")
						aiResponse, err := p.aiContextService.RespondWithContext(valueStr, contexto)
						if err == nil && aiResponse != "" {
							finalErrorMessage = aiResponse
							p.logger.Info(fmt.Sprintf("✅ Respuesta generada por IA: %s", aiResponse))
						} else {
							p.logger.Warn(fmt.Sprintf("⚠️ Error generando respuesta con IA: %v, usando mensaje por defecto", err))
							finalErrorMessage = "La respuesta no es válida. Por favor intenta de nuevo."
						}
					} else {
						// No hay contexto (flow no tiene idcontext) o no hay servicio de IA
						if contexto == nil {
							p.logger.Info("ℹ️ El flow no tiene idcontext, no se usará IA con contexto")
						}
						if p.aiContextService == nil {
							p.logger.Warn("⚠️ Servicio de IA con contexto no disponible")
						}
						finalErrorMessage = "La respuesta no es válida. Por favor intenta de nuevo."
					}
				}

				// Enviar mensaje de error al usuario
				if p.messagingService != nil {
					phone := session.ConversationID
					if idx := len(phone) - 1; idx > 0 {
						// Extraer solo el número de teléfono (antes del @)
						for i, ch := range phone {
							if ch == '@' {
								phone = phone[:i]
								break
							}
						}
					}

					errorMsg := &entities.Message{
						TenantID:       session.TenantID,
						InstanceID:     session.InstanceID,
						ConversationID: session.ConversationID,
						To:             phone,
						Direction:      "out",
						MessageData: entities.MessageData{
							Type: "text",
							Text: &entities.TextContent{
								Body: finalErrorMessage,
							},
						},
					}

					if err := p.messagingService.SendMessage(ctx, errorMsg); err != nil {
						p.logger.Error(fmt.Sprintf("Error enviando mensaje de error: %v", err))
					} else {
						p.logger.Info("📤 Mensaje de error enviado al usuario")
					}
				}

				// Marcar que sigue esperando respuesta (no avanzar)
				// Limpiar el flag de validación para permitir nueva validación
				session.SetVariable(variableName+"_validated", false)
				session.SetWaitingForResponse(variableName)

				// Si hay un errorMessage personalizado, ya se envió arriba
				// Retornar para que NO avance al siguiente nodo (esperar nueva respuesta)
				return &ProcessResult{
					WaitingForResponse: true,
					WaitingForVariable: variableName,
					StopFlow:           false,
					ErrorMessage:       finalErrorMessage,
				}, nil
			}

			p.logger.Info(fmt.Sprintf("✅ Validación IA exitosa: %s", reason))
		}
	}

	// Validar con reglas tradicionales
	if validationConfig != nil {
		validationError := p.validate(value, validationConfig)
		if validationError != "" {
			p.logger.Warn(fmt.Sprintf("Validation failed: %s", validationError))
			return &ProcessResult{
				StopFlow:     false,
				ErrorMessage: validationError,
			}, nil
		}
	}

	p.logger.Info(fmt.Sprintf("Validation passed for variable: %s", variableName))

	// Marcar que esta variable ya fue validada exitosamente para evitar bucles
	// Guardar un flag indicando que esta variable ya pasó la validación
	session.SetVariable(variableName+"_validated", true)

	// No espera respuesta, continúa automáticamente
	return &ProcessResult{
		WaitingForResponse: false,
		StopFlow:           false,
	}, nil
}

// validate valida un valor según las reglas
func (p *ResponseNodeProcessor) validate(value interface{}, rules map[string]interface{}) string {
	valueStr := fmt.Sprintf("%v", value)

	// Validar required
	if required, ok := rules["required"].(bool); ok && required {
		if valueStr == "" {
			return "El valor es requerido"
		}
	}

	// Validar minLength
	if minLength, ok := rules["minLength"].(float64); ok {
		if len(valueStr) < int(minLength) {
			return fmt.Sprintf("El valor debe tener al menos %d caracteres", int(minLength))
		}
	}

	// Validar maxLength
	if maxLength, ok := rules["maxLength"].(float64); ok {
		if len(valueStr) > int(maxLength) {
			return fmt.Sprintf("El valor no debe exceder %d caracteres", int(maxLength))
		}
	}

	// Validar pattern
	if pattern, ok := rules["pattern"].(string); ok && pattern != "" {
		matched, err := regexp.MatchString(pattern, valueStr)
		if err != nil {
			return fmt.Sprintf("Error en patrón de validación: %v", err)
		}
		if !matched {
			return "El valor no cumple con el formato esperado"
		}
	}

	return ""
}

// getContextoForSession obtiene el contexto para una sesión
// Primero busca el idcontext en el flow. Si el flow tiene idcontext, usa ese contexto.
// Si el flow no tiene idcontext, retorna nil (no se usa contexto)
func (p *ResponseNodeProcessor) getContextoForSession(ctx context.Context, session *entities.FlowSession) (*entities.Contexto, error) {
	if p.contextoRepo == nil {
		return nil, fmt.Errorf("contexto repository no configurado")
	}

	// 1. Buscar el flow de la sesión
	if p.flowRepo == nil {
		return nil, fmt.Errorf("flow repository no configurado")
	}

	flow, err := p.flowRepo.FindByID(ctx, session.FlowID)
	if err != nil {
		p.logger.Warn(fmt.Sprintf("⚠️ No se pudo encontrar el flow %s: %v", session.FlowID, err))
		return nil, fmt.Errorf("no se pudo encontrar el flow: %w", err)
	}

	// 2. Si el flow tiene idcontext, buscar ese contexto
	if flow.ContextID != "" {
		p.logger.Info(fmt.Sprintf("🔍 Buscando contexto con ID: %s (del flow %s)", flow.ContextID, session.FlowID))
		contexto, err := p.contextoRepo.FindByID(ctx, flow.ContextID)
		if err != nil {
			p.logger.Warn(fmt.Sprintf("⚠️ No se pudo encontrar el contexto %s: %v", flow.ContextID, err))
			return nil, fmt.Errorf("no se pudo encontrar el contexto %s: %w", flow.ContextID, err)
		}
		if contexto != nil {
			if !contexto.IsActive {
				p.logger.Warn(fmt.Sprintf("⚠️ El contexto %s existe pero no está activo", flow.ContextID))
				return nil, fmt.Errorf("el contexto %s no está activo", flow.ContextID)
			}
			p.logger.Info(fmt.Sprintf("✅ Contexto encontrado: %s", flow.ContextID))
			return contexto, nil
		}
		return nil, fmt.Errorf("contexto %s no encontrado", flow.ContextID)
	}

	// 3. Si el flow NO tiene idcontext, retornar nil (no usar contexto)
	p.logger.Info(fmt.Sprintf("ℹ️ El flow %s no tiene idcontext, no se usará contexto", session.FlowID))
	return nil, nil
}
