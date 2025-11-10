package flow

import (
	"context"
	"fmt"
	"reflect"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// ConditionNodeProcessor procesa nodos de tipo CONDITION
type ConditionNodeProcessor struct {
	logger           ports.Logger
	variableReplacer *VariableReplacer
}

// NewConditionNodeProcessor crea un nuevo procesador de condiciones
func NewConditionNodeProcessor(
	logger ports.Logger,
	variableReplacer *VariableReplacer,
) *ConditionNodeProcessor {
	return &ConditionNodeProcessor{
		logger:           logger,
		variableReplacer: variableReplacer,
	}
}

func (p *ConditionNodeProcessor) Process(ctx context.Context, session *entities.FlowSession, node *entities.FlowNode) (*ProcessResult, error) {
	p.logger.Info(fmt.Sprintf("Processing CONDITION node: %s", node.ID))

	// Extraer configuración
	config := node.Config
	conditionsConfig, _ := config["conditions"].([]interface{})

	// Evaluar cada condición
	conditionMet := false
	for _, condConfig := range conditionsConfig {
		condMap, ok := condConfig.(map[string]interface{})
		if !ok {
			continue
		}

		field, _ := condMap["field"].(string)
		operator, _ := condMap["operator"].(string)
		expectedValue := condMap["value"]

		// Obtener valor de la variable
		actualValue, exists := p.variableReplacer.GetNestedValue(field, session.Variables)
		if !exists {
			p.logger.Warn(fmt.Sprintf("Variable %s not found in session", field))
			continue
		}

		// Evaluar condición
		if p.evaluateCondition(actualValue, operator, expectedValue) {
			conditionMet = true
			break
		}
	}

	// Determinar el edge a seguir
	condition := "no"
	if conditionMet {
		condition = "yes"
	}

	p.logger.Info(fmt.Sprintf("Condition evaluated to: %s", condition))

	// El FlowEngine se encargará de buscar el edge correcto
	return &ProcessResult{
		WaitingForResponse: false,
		NextNodeID:         condition, // Usamos esto para indicar qué condición seguir
		StopFlow:           false,
	}, nil
}

// evaluateCondition evalúa una condición
func (p *ConditionNodeProcessor) evaluateCondition(actual interface{}, operator string, expected interface{}) bool {
	switch operator {
	case "equals", "==":
		return p.compareValues(actual, expected) == 0
	case "not_equals", "!=":
		return p.compareValues(actual, expected) != 0
	case "greater_than", ">":
		return p.compareValues(actual, expected) > 0
	case "less_than", "<":
		return p.compareValues(actual, expected) < 0
	case "greater_or_equal", ">=":
		return p.compareValues(actual, expected) >= 0
	case "less_or_equal", "<=":
		return p.compareValues(actual, expected) <= 0
	case "contains":
		return p.contains(actual, expected)
	default:
		p.logger.Warn(fmt.Sprintf("Unknown operator: %s", operator))
		return false
	}
}

// compareValues compara dos valores
func (p *ConditionNodeProcessor) compareValues(a, b interface{}) int {
	// Convertir a tipos comparables
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)

	// Comparar por tipo
	if reflect.TypeOf(a) == reflect.TypeOf(b) {
		switch a.(type) {
		case int, int64:
			aInt := a.(int64)
			bInt := b.(int64)
			if aInt < bInt {
				return -1
			} else if aInt > bInt {
				return 1
			}
			return 0
		case float64:
			aFloat := a.(float64)
			bFloat := b.(float64)
			if aFloat < bFloat {
				return -1
			} else if aFloat > bFloat {
				return 1
			}
			return 0
		case bool:
			aBool := a.(bool)
			bBool := b.(bool)
			if aBool == bBool {
				return 0
			}
			return -1
		}
	}

	// Comparación de strings
	if aStr < bStr {
		return -1
	} else if aStr > bStr {
		return 1
	}
	return 0
}

// contains verifica si un valor contiene otro
func (p *ConditionNodeProcessor) contains(haystack, needle interface{}) bool {
	haystackStr := fmt.Sprintf("%v", haystack)
	needleStr := fmt.Sprintf("%v", needle)
	return len(haystackStr) > 0 && len(needleStr) > 0 && 
		   (haystackStr == needleStr || len(haystackStr) > len(needleStr))
}


