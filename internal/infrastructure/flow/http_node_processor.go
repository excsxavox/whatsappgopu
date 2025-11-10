package flow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// HttpNodeProcessor procesa nodos de tipo HTTP
type HttpNodeProcessor struct {
	logger           ports.Logger
	variableReplacer *VariableReplacer
	httpClient       *http.Client
}

// NewHttpNodeProcessor crea un nuevo procesador HTTP
func NewHttpNodeProcessor(
	logger ports.Logger,
	variableReplacer *VariableReplacer,
) *HttpNodeProcessor {
	return &HttpNodeProcessor{
		logger:           logger,
		variableReplacer: variableReplacer,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *HttpNodeProcessor) Process(ctx context.Context, session *entities.FlowSession, node *entities.FlowNode) (*ProcessResult, error) {
	p.logger.Info(fmt.Sprintf("Processing HTTP node: %s", node.ID))

	// Extraer configuración
	config := node.Config
	method, _ := config["method"].(string)
	url, _ := config["url"].(string)
	headersConfig, _ := config["headers"].(map[string]interface{})
	bodyConfig, _ := config["body"].(map[string]interface{})
	responseVariable, _ := config["responseVariable"].(string)

	// Reemplazar variables en URL
	url = p.variableReplacer.ReplaceInString(url, session.Variables)

	// Reemplazar variables en headers
	headers := p.variableReplacer.ReplaceInMap(headersConfig, session.Variables)

	// Reemplazar variables en body
	body := p.variableReplacer.ReplaceInMap(bodyConfig, session.Variables)

	// Preparar body JSON
	var bodyBytes []byte
	var err error
	if len(body) > 0 {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			p.logger.Error(fmt.Sprintf("Error marshaling HTTP body: %v", err))
			// Guardar error en variable
			session.SetVariable(responseVariable, map[string]interface{}{
				"error":   true,
				"message": fmt.Sprintf("Error preparing request: %v", err),
			})
			return &ProcessResult{StopFlow: false}, nil
		}
	}

	// Crear request
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		p.logger.Error(fmt.Sprintf("Error creating HTTP request: %v", err))
		session.SetVariable(responseVariable, map[string]interface{}{
			"error":   true,
			"message": fmt.Sprintf("Error creating request: %v", err),
		})
		return &ProcessResult{StopFlow: false}, nil
	}

	// Aplicar headers
	for key, value := range headers {
		if strValue, ok := value.(string); ok {
			req.Header.Set(key, strValue)
		}
	}

	// Ejecutar request
	p.logger.Info(fmt.Sprintf("Executing HTTP %s to %s", method, url))
	resp, err := p.httpClient.Do(req)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Error executing HTTP request: %v", err))
		session.SetVariable(responseVariable, map[string]interface{}{
			"error":   true,
			"message": fmt.Sprintf("Request failed: %v", err),
		})
		return &ProcessResult{StopFlow: false}, nil
	}
	defer resp.Body.Close()

	// Leer respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Error reading HTTP response: %v", err))
		session.SetVariable(responseVariable, map[string]interface{}{
			"error":   true,
			"message": fmt.Sprintf("Error reading response: %v", err),
		})
		return &ProcessResult{StopFlow: false}, nil
	}

	// Parsear JSON de respuesta
	var responseData interface{}
	if err := json.Unmarshal(respBody, &responseData); err != nil {
		// Si no es JSON, guardar como string
		responseData = string(respBody)
	}

	// Guardar respuesta en variable
	responseMap := map[string]interface{}{
		"status_code": resp.StatusCode,
		"data":        responseData,
		"error":       resp.StatusCode >= 400,
	}

	session.SetVariable(responseVariable, responseMap)

	p.logger.Info(fmt.Sprintf("HTTP request completed with status %d", resp.StatusCode))

	// No espera respuesta del usuario, continúa automáticamente
	return &ProcessResult{
		WaitingForResponse: false,
		StopFlow:           false,
	}, nil
}


