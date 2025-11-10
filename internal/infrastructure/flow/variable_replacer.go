package flow

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// VariableReplacer maneja el reemplazo de variables en strings
type VariableReplacer struct{}

// NewVariableReplacer crea un nuevo replacer
func NewVariableReplacer() *VariableReplacer {
	return &VariableReplacer{}
}

// ReplaceInString reemplaza variables en un string
// Soporta: {variable} y [variable]
func (vr *VariableReplacer) ReplaceInString(text string, variables map[string]interface{}) string {
	if text == "" {
		return text
	}

	// Patrón para {variable}
	re1 := regexp.MustCompile(`\{([a-zA-Z0-9_\.]+)\}`)
	text = re1.ReplaceAllStringFunc(text, func(match string) string {
		// Extraer nombre de variable (sin { })
		varName := match[1 : len(match)-1]
		return vr.getVariableValue(varName, variables)
	})

	// Patrón para [variable]
	re2 := regexp.MustCompile(`\[([a-zA-Z0-9_\.]+)\]`)
	text = re2.ReplaceAllStringFunc(text, func(match string) string {
		// Extraer nombre de variable (sin [ ])
		varName := match[1 : len(match)-1]
		return vr.getVariableValue(varName, variables)
	})

	return text
}

// ReplaceInMap reemplaza variables en un map
func (vr *VariableReplacer) ReplaceInMap(data map[string]interface{}, variables map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		result[key] = vr.replaceInValue(value, variables)
	}

	return result
}

// replaceInValue reemplaza variables en cualquier tipo de valor
func (vr *VariableReplacer) replaceInValue(value interface{}, variables map[string]interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return vr.ReplaceInString(v, variables)
	case map[string]interface{}:
		return vr.ReplaceInMap(v, variables)
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = vr.replaceInValue(item, variables)
		}
		return result
	default:
		return value
	}
}

// getVariableValue obtiene el valor de una variable (soporta notación de punto)
func (vr *VariableReplacer) getVariableValue(varName string, variables map[string]interface{}) string {
	// Si la variable contiene punto (ej: response.valid), navegar el objeto
	if strings.Contains(varName, ".") {
		parts := strings.Split(varName, ".")
		current := variables

		for i, part := range parts {
			if val, exists := current[part]; exists {
				// Si es el último elemento, devolver el valor
				if i == len(parts)-1 {
					return vr.valueToString(val)
				}
				// Si no es el último, verificar que sea un map
				if nextMap, ok := val.(map[string]interface{}); ok {
					current = nextMap
				} else {
					return fmt.Sprintf("{%s}", varName) // No se pudo navegar
				}
			} else {
				return fmt.Sprintf("{%s}", varName) // Variable no existe
			}
		}
	}

	// Variable simple (sin punto)
	if val, exists := variables[varName]; exists {
		return vr.valueToString(val)
	}

	// Variable no existe, devolver el placeholder original
	return fmt.Sprintf("{%s}", varName)
}

// valueToString convierte cualquier valor a string
func (vr *VariableReplacer) valueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int, int64, float64, bool:
		return fmt.Sprintf("%v", v)
	case map[string]interface{}, []interface{}:
		// Para objetos/arrays, convertir a JSON
		bytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(bytes)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetNestedValue obtiene un valor anidado usando notación de punto
func (vr *VariableReplacer) GetNestedValue(varName string, variables map[string]interface{}) (interface{}, bool) {
	if !strings.Contains(varName, ".") {
		val, exists := variables[varName]
		return val, exists
	}

	parts := strings.Split(varName, ".")
	current := variables

	for i, part := range parts {
		if val, exists := current[part]; exists {
			if i == len(parts)-1 {
				return val, true
			}
			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return nil, false
			}
		} else {
			return nil, false
		}
	}

	return nil, false
}


