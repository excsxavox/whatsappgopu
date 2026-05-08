package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"
)

// AIContextService servicio para responder preguntas usando el contexto de la empresa
type AIContextService struct {
	apiKey     string
	httpClient *http.Client
	logger     ports.Logger
}

type openAIChatRequest struct {
	Model    string              `json:"model"`
	Messages []openAIChatMessage `json:"messages"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// NewAIContextService crea un nuevo servicio de IA con contexto
func NewAIContextService(apiKey string, logger ports.Logger) *AIContextService {
	return &AIContextService{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		logger:     logger,
	}
}

// RespondWithContext genera una respuesta usando el contexto de la empresa
// Retorna la respuesta generada por la IA basada en el contexto
func (s *AIContextService) RespondWithContext(userQuestion string, contexto *entities.Contexto) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY no configurada")
	}

	if contexto == nil {
		return "", fmt.Errorf("contexto es nil")
	}

	s.logger.Info(fmt.Sprintf("🤖 Generando respuesta con contexto para pregunta: '%s'", userQuestion))

	// Construir el prompt completo usando el contexto
	fullPrompt := contexto.GetFullPrompt()

	// Construir el mensaje del sistema con el contexto
	systemPrompt := fullPrompt + "\n\nResponde de manera amigable y profesional. Si no tienes la información necesaria, sé honesto y pide más detalles."

	// Preparar la solicitud
	reqBody := openAIChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openAIChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userQuestion},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error calling OpenAI: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		s.logger.Error(fmt.Sprintf("OpenAI API error (status %d): %s", resp.StatusCode, string(body)))
		return "", fmt.Errorf("openAI API returned status %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp openAIChatResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("openAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	response := openAIResp.Choices[0].Message.Content
	s.logger.Info(fmt.Sprintf("✅ Respuesta generada: %s", response))

	return response, nil
}

// IsQuestion detecta si el mensaje del usuario es una pregunta
// Retorna true si parece ser una pregunta
func (s *AIContextService) IsQuestion(userMessage string) bool {
	// Detectar preguntas comunes
	questionWords := []string{"qué", "cuál", "cuándo", "dónde", "cómo", "por qué", "quién", "cuánto", "cuánta", "cuántos", "cuántas"}

	// Verificar si tiene signos de interrogación
	if len(userMessage) > 0 && (userMessage[0] == '¿' || userMessage[len(userMessage)-1] == '?') {
		return true
	}

	// Buscar palabras de pregunta en el mensaje
	userMessageLower := fmt.Sprintf(" %s ", userMessage) // Agregar espacios para mejor matching
	for _, word := range questionWords {
		if len(userMessage) >= len(word) {
			wordWithSpaces := fmt.Sprintf(" %s ", word)
			for i := 0; i <= len(userMessageLower)-len(wordWithSpaces); i++ {
				if userMessageLower[i:i+len(wordWithSpaces)] == wordWithSpaces {
					return true
				}
			}
		}
	}

	return false
}
