package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"whatsapp-api-go/internal/domain/ports"

	"golang.org/x/oauth2/google"
)

// GoogleTTSService servicio para convertir texto a voz usando Google Cloud Text-to-Speech API
// Soporta tanto el modelo tradicional como el nuevo Gemini Pro TTS
type GoogleTTSService struct {
	apiKey         string // API key (para modelo standard) o JSON de cuenta de servicio (para Gemini)
	serviceAccount bool   // true si usa autenticación con cuenta de servicio
	logger         ports.Logger
	httpClient     *http.Client
	voiceName      string // Nombre de la voz (ej: "Erinome" para Gemini Pro TTS, "es-ES-Neural2-A" para modelo tradicional)
	model          string // Modelo a usar: "gemini-2.5-pro-tts" o "standard"
	language       string // Código de idioma: "es-ES" (España) o "es-MX" (México)
	prompt         string // Prompt opcional para Gemini Pro TTS (ej: "Read aloud in a warm, welcoming tone.")
}

// NewGoogleTTSService crea un nuevo servicio TTS de Google
// apiKeyOrJSON: Clave de API de Google Cloud (para modelo standard) o JSON de cuenta de servicio (para Gemini)
// voiceName: Nombre de la voz (por defecto "es-ES-Neural2-A" para modelo standard)
// model: Modelo a usar - "gemini" para Gemini Pro TTS o "standard" para modelo tradicional
func NewGoogleTTSService(apiKeyOrJSON string, logger ports.Logger, voiceName string, model string, language string, prompt string) *GoogleTTSService {
	if model == "" {
		model = "standard" // Por defecto usar modelo standard (no requiere Vertex AI, funciona con API key)
	}
	if voiceName == "" {
		if model == "gemini" {
			voiceName = "Erinome" // Voz por defecto para Gemini Pro TTS (femenina, español)
		} else {
			voiceName = "es-ES-Neural2-A" // Voz por defecto para modelo standard (femenina, español de España)
		}
	}
	if language == "" {
		language = "es-ES" // Español de España (castellano) por defecto
	}
	if prompt == "" && model == "gemini" {
		prompt = "Read aloud in a warm, welcoming tone." // Prompt por defecto para Gemini
	}

	// Detectar si es JSON de cuenta de servicio (empieza con {) o API key
	serviceAccount := strings.TrimSpace(apiKeyOrJSON) != "" && strings.HasPrefix(strings.TrimSpace(apiKeyOrJSON), "{")

	return &GoogleTTSService{
		apiKey:         apiKeyOrJSON,
		serviceAccount: serviceAccount,
		logger:         logger,
		httpClient:     &http.Client{},
		voiceName:      voiceName,
		model:          model,
		language:       language,
		prompt:         prompt,
	}
}

// GoogleTTSRequest estructura para la petición a Google Cloud TTS API (modelo tradicional)
type GoogleTTSRequest struct {
	Input struct {
		Text string `json:"text"`
	} `json:"input"`
	Voice struct {
		LanguageCode string `json:"languageCode"`
		Name         string `json:"name"`
		SSMLGender   string `json:"ssmlGender"`
	} `json:"voice"`
	AudioConfig struct {
		AudioEncoding string  `json:"audioEncoding"`
		SpeakingRate  float64 `json:"speakingRate,omitempty"`
		Pitch         float64 `json:"pitch,omitempty"`
	} `json:"audioConfig"`
}

// GoogleGeminiTTSRequest estructura para la petición a Gemini Pro TTS API (v1beta1)
type GoogleGeminiTTSRequest struct {
	Input struct {
		Text   string `json:"text"`
		Prompt string `json:"prompt,omitempty"` // Prompt opcional para controlar el tono
	} `json:"input"`
	Voice struct {
		LanguageCode string `json:"languageCode"`
		ModelName    string `json:"modelName"` // "gemini-2.5-pro-tts"
		Name         string `json:"name"`      // Nombre de la voz (ej: "Erinome")
	} `json:"voice"`
	AudioConfig struct {
		AudioEncoding string  `json:"audioEncoding"` // "MP3" o "LINEAR16"
		Pitch         float64 `json:"pitch,omitempty"`
		SpeakingRate  float64 `json:"speakingRate,omitempty"`
	} `json:"audioConfig"`
}

// GoogleTTSResponse estructura para la respuesta de Google Cloud TTS API
type GoogleTTSResponse struct {
	AudioContent string `json:"audioContent"`
}

// TextToSpeech convierte texto a audio usando Google Cloud TTS API
// Retorna los bytes del audio en formato MP3
// Soporta pausas usando [Pausa] o [Pausa: X] donde X es el número de segundos
func (s *GoogleTTSService) TextToSpeech(text string) ([]byte, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_TTS_API_KEY no configurada")
	}

	if text == "" {
		return nil, fmt.Errorf("texto vacío")
	}

	// Procesar pausas en el texto
	text = s.processPauses(text)

	// Preprocesar texto para mejorar pronunciación en español de España (castellano)
	text = s.preprocessForCastilianSpanish(text)

	s.logger.Info(fmt.Sprintf("🎤 [Google TTS] Generando audio con modelo %s, voz %s, idioma %s para texto: %s", s.model, s.voiceName, s.language, text))

	var url string
	var jsonData []byte
	var err error

	if s.model == "gemini" {
		// Usar Gemini Pro TTS (v1beta1)
		if s.serviceAccount {
			// Para cuenta de servicio, no usar ?key= en la URL
			url = "https://texttospeech.googleapis.com/v1beta1/text:synthesize"
		} else {
			// Para API key, usar ?key= en la URL
			url = fmt.Sprintf("https://texttospeech.googleapis.com/v1beta1/text:synthesize?key=%s", s.apiKey)
		}

		reqBody := GoogleGeminiTTSRequest{}
		reqBody.Input.Text = text
		reqBody.Input.Prompt = s.prompt // Prompt para controlar el tono
		reqBody.Voice.LanguageCode = s.language
		reqBody.Voice.ModelName = "gemini-2.5-pro-tts"
		reqBody.Voice.Name = s.voiceName
		reqBody.AudioConfig.AudioEncoding = "MP3" // Para WhatsApp usamos MP3
		reqBody.AudioConfig.Pitch = 0.0
		reqBody.AudioConfig.SpeakingRate = 1.0

		jsonData, err = json.Marshal(reqBody)
	} else {
		// Usar modelo tradicional (v1)
		if s.serviceAccount {
			url = "https://texttospeech.googleapis.com/v1/text:synthesize"
		} else {
			url = fmt.Sprintf("https://texttospeech.googleapis.com/v1/text:synthesize?key=%s", s.apiKey)
		}

		reqBody := GoogleTTSRequest{}
		reqBody.Input.Text = text
		reqBody.Voice.LanguageCode = s.language
		reqBody.Voice.Name = s.voiceName
		reqBody.Voice.SSMLGender = "FEMALE" // o "MALE" según la voz
		reqBody.AudioConfig.AudioEncoding = "MP3"
		reqBody.AudioConfig.SpeakingRate = 1.0
		reqBody.AudioConfig.Pitch = 0.0

		jsonData, err = json.Marshal(reqBody)
	}
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Crear request
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Si usa cuenta de servicio, obtener token OAuth2
	if s.serviceAccount {
		token, err := s.getAccessToken()
		if err != nil {
			return nil, fmt.Errorf("error obteniendo token de acceso: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	// Ejecutar request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		s.logger.Error(fmt.Sprintf("Google TTS API error (status %d): %s", resp.StatusCode, string(body)))
		return nil, fmt.Errorf("google TTS API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parsear respuesta
	var ttsResponse GoogleTTSResponse
	if err := json.Unmarshal(body, &ttsResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Decodificar audio base64
	audioData, err := s.decodeBase64Audio(ttsResponse.AudioContent)
	if err != nil {
		return nil, fmt.Errorf("error decoding audio: %w", err)
	}

	s.logger.Info(fmt.Sprintf("✅ [Google TTS] Audio generado: %d bytes", len(audioData)))
	return audioData, nil
}

// getAccessToken obtiene un token de acceso OAuth2 desde el JSON de la cuenta de servicio
func (s *GoogleTTSService) getAccessToken() (string, error) {
	ctx := context.Background()

	// Crear credenciales desde el JSON
	creds, err := google.CredentialsFromJSON(ctx, []byte(s.apiKey), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return "", fmt.Errorf("error cargando credenciales: %w", err)
	}

	// Obtener token
	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("error obteniendo token: %w", err)
	}

	return token.AccessToken, nil
}

// decodeBase64Audio decodifica el audio en base64 que retorna Google TTS
func (s *GoogleTTSService) decodeBase64Audio(base64Audio string) ([]byte, error) {
	// Google TTS retorna el audio en base64
	// Necesitamos decodificarlo
	return base64.StdEncoding.DecodeString(base64Audio)
}

// processPauses procesa las etiquetas de pausa en el texto
func (s *GoogleTTSService) processPauses(text string) string {
	re := regexp.MustCompile(`\[Pausa(?::\s*(\d+))?\]`)

	result := re.ReplaceAllStringFunc(text, func(match string) string {
		parts := re.FindStringSubmatch(match)

		if len(parts) > 1 && parts[1] != "" {
			seconds := parts[1]
			if seconds == "1" {
				return ". "
			} else if seconds == "2" {
				return "... "
			} else {
				return ".... "
			}
		}

		return ", "
	})

	return result
}

// preprocessForCastilianSpanish preprocesa el texto para mejorar la pronunciación en español de España
func (s *GoogleTTSService) preprocessForCastilianSpanish(text string) string {
	// Normalizar espacios y puntuación
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = regexp.MustCompile(`\s*,\s*`).ReplaceAllString(text, ", ")
	text = regexp.MustCompile(`\s*\.\s*`).ReplaceAllString(text, ". ")

	// Asegurar que las preguntas tengan la entonación correcta
	text = regexp.MustCompile(`\?+`).ReplaceAllString(text, "?")
	text = regexp.MustCompile(`\!+`).ReplaceAllString(text, "!")

	// Reemplazar palabras americanas por castellanas
	replacements := map[string]string{
		"celular":     "móvil",
		"computadora": "ordenador",
		"carro":       "coche",
		"apartamento": "piso",
	}

	for american, castilian := range replacements {
		pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(american) + `\b`)
		text = pattern.ReplaceAllString(text, castilian)
	}

	return text
}
