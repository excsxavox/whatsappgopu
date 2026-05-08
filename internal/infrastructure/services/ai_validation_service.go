package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"whatsapp-api-go/internal/domain/ports"
)

type AIValidationService struct {
	apiKey     string
	httpClient *http.Client
	logger     ports.Logger
}

type openAIValidationRequest struct {
	Model      string                    `json:"model"`
	Messages   []openAIValidationMessage `json:"messages"`
	Tools      []openAITool              `json:"tools"`
	ToolChoice map[string]interface{}    `json:"tool_choice"`
}

type openAIValidationMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAITool struct {
	Type     string             `json:"type"`
	Function openAIToolFunction `json:"function"`
}

type openAIToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type openAIValidationResponse struct {
	Choices []struct {
		Message struct {
			ToolCalls []struct {
				Function struct {
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
}

type validationResult struct {
	IsValid bool   `json:"is_valid"`
	Reason  string `json:"reason,omitempty"`
}

func NewAIValidationService(apiKey string, logger ports.Logger) *AIValidationService {
	return &AIValidationService{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		logger:     logger,
	}
}

// ValidateWithAI valida la respuesta del usuario usando OpenAI
func (s *AIValidationService) ValidateWithAI(userResponse, validationPrompt string) (bool, string, error) {
	if s.apiKey == "" {
		return true, "", nil // Si no hay API key, no validar
	}

	s.logger.Info(fmt.Sprintf("🤖 Validando respuesta con IA: '%s'", userResponse))

	// Construir el prompt completo
	systemPrompt := "Eres un asistente que valida respuestas de usuarios en un chatbot de WhatsApp. Responde SOLO si la respuesta es válida o no según las instrucciones."
	userPrompt := fmt.Sprintf(`Instrucciones de validación: %s

Respuesta del usuario: "%s"

¿Es válida esta respuesta según las instrucciones?`, validationPrompt, userResponse)

	// Preparar la solicitud con Function Calling
	reqBody := openAIValidationRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openAIValidationMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Tools: []openAITool{
			{
				Type: "function",
				Function: openAIToolFunction{
					Name:        "validate_response",
					Description: "Valida si la respuesta del usuario cumple con los criterios especificados",
					Parameters: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"is_valid": map[string]interface{}{
								"type":        "boolean",
								"description": "true si la respuesta es válida, false si no lo es",
							},
							"reason": map[string]interface{}{
								"type":        "string",
								"description": "Razón por la cual la respuesta es válida o inválida",
							},
						},
						"required": []string{"is_valid"},
					},
				},
			},
		},
		ToolChoice: map[string]interface{}{
			"type": "function",
			"function": map[string]string{
				"name": "validate_response",
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return false, "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("error calling OpenAI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, "", fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	var openAIResp openAIValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return false, "", fmt.Errorf("error decoding response: %w", err)
	}

	// Extraer el resultado de la función
	if len(openAIResp.Choices) == 0 || len(openAIResp.Choices[0].Message.ToolCalls) == 0 {
		return false, "", fmt.Errorf("no function call in response")
	}

	var result validationResult
	if err := json.Unmarshal([]byte(openAIResp.Choices[0].Message.ToolCalls[0].Function.Arguments), &result); err != nil {
		return false, "", fmt.Errorf("error parsing validation result: %w", err)
	}

	if result.IsValid {
		s.logger.Info(fmt.Sprintf("✅ Validación IA: VÁLIDA - %s", result.Reason))
	} else {
		s.logger.Info(fmt.Sprintf("❌ Validación IA: INVÁLIDA - %s", result.Reason))
	}

	return result.IsValid, result.Reason, nil
}
