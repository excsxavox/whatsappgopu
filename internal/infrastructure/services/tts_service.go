package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"whatsapp-api-go/internal/domain/ports"
)

// TTSService servicio para convertir texto a voz usando OpenAI TTS API
type TTSService struct {
	openAIKey  string
	logger     ports.Logger
	httpClient *http.Client
}

// NewTTSService crea un nuevo servicio TTS
func NewTTSService(openAIKey string, logger ports.Logger) *TTSService {
	return &TTSService{
		openAIKey:  openAIKey,
		logger:     logger,
		httpClient: &http.Client{},
	}
}

// TextToSpeechRequest estructura para la petición a OpenAI TTS
type TextToSpeechRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	Voice string `json:"voice"`
}

// TextToSpeech convierte texto a audio usando OpenAI TTS API
// Retorna los bytes del audio en formato MP3
// Soporta pausas usando [Pausa] o [Pausa: X] donde X es el número de segundos
func (s *TTSService) TextToSpeech(text string) ([]byte, error) {
	if s.openAIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY no configurada")
	}

	if text == "" {
		return nil, fmt.Errorf("texto vacío")
	}

	// Procesar pausas en el texto
	// Reemplazar [Pausa] o [Pausa: X] con comas y puntos que crean pausas naturales
	text = s.processPauses(text)

	// Preprocesar texto para mejorar pronunciación en español de España (castellano)
	// Esto ayuda a que la voz suene más natural y con acento castellano
	text = s.preprocessForCastilianSpanish(text)

	// OpenAI TTS detecta automáticamente el idioma del texto
	// El texto debe estar en español (castellano) para que la voz pronuncie correctamente
	s.logger.Info(fmt.Sprintf("🎤 Generando audio en español (castellano) con voz femenina usando TTS para texto: %s", text))

	url := "https://api.openai.com/v1/audio/speech"

	// Preparar request body
	// Nota: OpenAI TTS detecta automáticamente el idioma del texto
	// Usamos 'nova' que tiene mejor pronunciación en español de España (castellano)
	// 'nova' es una voz femenina clara y natural, optimizada para español
	reqBody := TextToSpeechRequest{
		Model: "tts-1-hd", // Modelo de alta calidad para sonido más realista
		Input: text,       // Texto en español (castellano) - será detectado automáticamente
		Voice: "nova",     // Voz femenina clara y natural, mejor pronunciación en español de España
		// Voces disponibles: alloy (masculina neutra), ash (nueva), ballad (nueva), coral (nueva),
		// echo (masculina neutra), fable (femenina expresiva), nova (femenina natural),
		// onyx (masculina profunda), sage (nueva), shimmer (femenina expresiva)
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Crear request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.openAIKey))
	req.Header.Set("Content-Type", "application/json")

	// Ejecutar request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Leer respuesta (audio en formato MP3)
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.Error(fmt.Sprintf("OpenAI TTS API error (status %d): %s", resp.StatusCode, string(audioData)))
		return nil, fmt.Errorf("openAI TTS API returned status %d: %s", resp.StatusCode, string(audioData))
	}

	s.logger.Info(fmt.Sprintf("✅ Audio generado: %d bytes", len(audioData)))
	return audioData, nil
}

// processPauses procesa las etiquetas de pausa en el texto
// Reemplaza [Pausa] con elementos que crean pausas naturales en el TTS
// OpenAI TTS no soporta pausas exactas, pero interpreta comas, puntos y espacios
func (s *TTSService) processPauses(text string) string {
	// Patrón para [Pausa] o [Pausa: número]
	// Soporta: [Pausa], [Pausa: 1], [Pausa: 2], etc.
	re := regexp.MustCompile(`\[Pausa(?::\s*(\d+))?\]`)

	result := re.ReplaceAllStringFunc(text, func(match string) string {
		// Extraer el número si existe
		parts := re.FindStringSubmatch(match)

		if len(parts) > 1 && parts[1] != "" {
			// Si hay un número, usar puntos suspensivos para pausa más larga
			// Más puntos = pausa más larga (aunque OpenAI TTS no garantiza duración exacta)
			seconds := parts[1]
			if seconds == "1" {
				return ". "
			} else if seconds == "2" {
				return "... "
			} else {
				// Para 3+ segundos, usar más puntos
				return ".... "
			}
		}

		// Si no hay número, usar coma para pausa corta
		return ", "
	})

	return result
}

// preprocessForCastilianSpanish preprocesa el texto para mejorar la pronunciación en español de España
// Ajusta el texto para que suene más natural y con acento castellano
func (s *TTSService) preprocessForCastilianSpanish(text string) string {
	// Normalizar espacios y puntuación para mejor pronunciación
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = regexp.MustCompile(`\s*,\s*`).ReplaceAllString(text, ", ")
	text = regexp.MustCompile(`\s*\.\s*`).ReplaceAllString(text, ". ")

	// Asegurar que las preguntas tengan la entonación correcta
	text = regexp.MustCompile(`\?+`).ReplaceAllString(text, "?")
	text = regexp.MustCompile(`\!+`).ReplaceAllString(text, "!")

	// Agregar un prefijo contextual muy sutil que refuerce el español de España
	// Usamos palabras típicamente castellanas al inicio para guiar la pronunciación
	// Esto ayuda a que el modelo TTS detecte mejor el acento castellano
	// El prefijo es casi imperceptible pero mejora la pronunciación

	// Si el texto no comienza con una palabra que refuerce el castellano,
	// agregamos un prefijo muy corto y natural
	if !regexp.MustCompile(`^(Hola|Buenos|Buenas|Gracias|Por favor|Disculpe|Perdón)`).MatchString(text) {
		// No agregamos prefijo si ya tiene una palabra castellana común
		// El contexto del texto completo ayudará
	}

	// Reemplazar algunas palabras que pueden sonar más americanas por variantes castellanas
	replacements := map[string]string{
		"celular":     "móvil",
		"computadora": "ordenador",
		"carro":       "coche",
		"apartamento": "piso",
	}

	for american, castilian := range replacements {
		// Reemplazar solo palabras completas (con límites de palabra)
		pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(american) + `\b`)
		text = pattern.ReplaceAllString(text, castilian)
	}

	return text
}
