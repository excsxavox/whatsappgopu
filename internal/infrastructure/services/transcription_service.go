package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"whatsapp-api-go/internal/domain/ports"
)

type TranscriptionService struct {
	openAIKey  string
	logger     ports.Logger
	httpClient *http.Client
}

func NewTranscriptionService(openAIKey string, logger ports.Logger) *TranscriptionService {
	return &TranscriptionService{
		openAIKey:  openAIKey,
		logger:     logger,
		httpClient: &http.Client{},
	}
}

// TranscribeAudio transcribe un archivo de audio a texto usando OpenAI Whisper
func (s *TranscriptionService) TranscribeAudio(audioData []byte, filename string) (string, error) {
	url := "https://api.openai.com/v1/audio/transcriptions"

	// Crear multipart request
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Agregar archivo de audio
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("error creating form file: %w", err)
	}

	if _, err := part.Write(audioData); err != nil {
		return "", fmt.Errorf("error writing audio data: %w", err)
	}

	// Agregar modelo
	writer.WriteField("model", "whisper-1")

	// Agregar idioma español (opcional, mejora precisión)
	writer.WriteField("language", "es")

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("error closing multipart writer: %w", err)
	}

	// Crear request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.openAIKey))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Ejecutar request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.Error(fmt.Sprintf("OpenAI API error: %s", string(respBody)))
		return "", fmt.Errorf("openAI API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parsear respuesta
	var result struct {
		Text string `json:"text"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	s.logger.Info(fmt.Sprintf("✅ Audio transcrito: %s", result.Text))
	return result.Text, nil
}
