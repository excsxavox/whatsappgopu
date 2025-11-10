package flow

import (
	"context"
	"fmt"
	"regexp"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// ResponseNodeProcessor procesa nodos de tipo RESPONSE
type ResponseNodeProcessor struct {
	logger           ports.Logger
	variableReplacer *VariableReplacer
}

// NewResponseNodeProcessor crea un nuevo procesador de respuestas
func NewResponseNodeProcessor(
	logger ports.Logger,
	variableReplacer *VariableReplacer,
) *ResponseNodeProcessor {
	return &ResponseNodeProcessor{
		logger:           logger,
		variableReplacer: variableReplacer,
	}
}

func (p *ResponseNodeProcessor) Process(ctx context.Context, session *entities.FlowSession, node *entities.FlowNode) (*ProcessResult, error) {
	p.logger.Info(fmt.Sprintf("Processing RESPONSE node: %s", node.ID))

	// Extraer configuración
	config := node.Config
	variableName, _ := config["variableName"].(string)
	validationConfig, _ := config["validation"].(map[string]interface{})

	// Obtener valor de la variable
	value, exists := session.GetVariable(variableName)
	if !exists {
		p.logger.Warn(fmt.Sprintf("Variable %s not found in session", variableName))
		return &ProcessResult{
			StopFlow:     false,
			ErrorMessage: fmt.Sprintf("Variable %s not found", variableName),
		}, nil
	}

	// Validar
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


