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

	// Buscar condición en formato "si" (nuevo formato del flow designer)
	var field string
	var operator string
	var expectedValue interface{}
	var conditionMet bool

	if siConfig, ok := config["si"].(map[string]interface{}); ok {
		// Formato nuevo: { "si": { "field": "...", "value": "...", "operator": "..." } }
		field, _ = siConfig["field"].(string)
		operator, _ = siConfig["operator"].(string)
		expectedValue = siConfig["value"]

		// Si no hay operador, usar "equals" por defecto
		if operator == "" {
			operator = "equals"
		}

		p.logger.Info(fmt.Sprintf("📊 Evaluating condition: %s %s %v", field, operator, expectedValue))

		// Obtener valor de la variable
		actualValue, exists := session.GetVariable(field)
		if !exists {
			p.logger.Warn(fmt.Sprintf("⚠️ Variable %s not found in session", field))
			conditionMet = false
		} else {
			p.logger.Info(fmt.Sprintf("📝 Variable %s = %v (type: %T)", field, actualValue, actualValue))
			// Evaluar condición
			conditionMet = p.evaluateCondition(actualValue, operator, expectedValue)
		}
	} else {
		// Formato antiguo: { "conditions": [ ... ] }
		conditionsConfig, _ := config["conditions"].([]interface{})

		// Evaluar cada condición
		for _, condConfig := range conditionsConfig {
			condMap, ok := condConfig.(map[string]interface{})
			if !ok {
				continue
			}

			field, _ = condMap["field"].(string)
			operator, _ = condMap["operator"].(string)
			expectedValue = condMap["value"]

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
	}

	// Determinar el edge a seguir
	condition := "no"
	if conditionMet {
		condition = "yes"
	}

	p.logger.Info(fmt.Sprintf("✅ Condition evaluated to: %s", condition))

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
	case "length_equals":
		return p.lengthEquals(actual, expected)
	case "has_digits":
		return p.hasDigits(actual, expected)
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

// lengthEquals verifica si la longitud del texto es igual al valor esperado
func (p *ConditionNodeProcessor) lengthEquals(actual interface{}, expected interface{}) bool {
	actualStr := fmt.Sprintf("%v", actual)

	// Convertir expected a int
	var expectedLen int
	switch v := expected.(type) {
	case int:
		expectedLen = v
	case int32:
		expectedLen = int(v)
	case int64:
		expectedLen = int(v)
	case float64:
		expectedLen = int(v)
	case string:
		fmt.Sscanf(v, "%d", &expectedLen)
	default:
		p.logger.Warn(fmt.Sprintf("Cannot convert expected value to int: %v (%T)", expected, expected))
		return false
	}

	actualLen := len(actualStr)
	p.logger.Info(fmt.Sprintf("📏 Length check: '%s' has %d chars, expected %d", actualStr, actualLen, expectedLen))

	return actualLen == expectedLen
}

// hasDigits verifica si el texto tiene exactamente N dígitos numéricos (0-9)
func (p *ConditionNodeProcessor) hasDigits(actual interface{}, expected interface{}) bool {
	actualStr := fmt.Sprintf("%v", actual)

	// Convertir expected a int
	var expectedDigits int
	switch v := expected.(type) {
	case int:
		expectedDigits = v
	case int32:
		expectedDigits = int(v)
	case int64:
		expectedDigits = int(v)
	case float64:
		expectedDigits = int(v)
	case string:
		fmt.Sscanf(v, "%d", &expectedDigits)
	default:
		p.logger.Warn(fmt.Sprintf("Cannot convert expected value to int: %v (%T)", expected, expected))
		return false
	}

	// Contar solo dígitos
	digitCount := 0
	for _, char := range actualStr {
		if char >= '0' && char <= '9' {
			digitCount++
		}
	}

	p.logger.Info(fmt.Sprintf("🔢 Digit check: '%s' has %d digits, expected %d", actualStr, digitCount, expectedDigits))

	return digitCount == expectedDigits
}
