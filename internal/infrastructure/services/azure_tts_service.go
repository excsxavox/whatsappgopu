package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"whatsapp-api-go/internal/domain/ports"
)

// AzureTTSService servicio para convertir texto a voz usando Azure Speech Service
type AzureTTSService struct {
	subscriptionKey string
	region          string
	logger          ports.Logger
	httpClient      *http.Client
	voiceName       string // Nombre de la voz (ej: "es-CO-GonzaloNeural" para español colombiano)
	language        string // Código de idioma: "es-CO" (Colombia), "es-MX" (México) o "es-ES" (España)
}

// NewAzureTTSService crea un nuevo servicio TTS de Azure
// subscriptionKey: Clave de suscripción de Azure Speech Service
// region: Región de Azure (ej: "eastus", "westus")
// voiceName: Nombre de la voz (por defecto "es-CO-GonzaloNeural" para español colombiano)
func NewAzureTTSService(subscriptionKey string, region string, logger ports.Logger, voiceName string, language string) *AzureTTSService {
	if voiceName == "" {
		voiceName = "es-CO-GonzaloNeural" // Voz por defecto (masculina, muy natural y expresiva, español colombiano)
	}
	if language == "" {
		language = "es-CO" // Español colombiano por defecto (muy natural para latinoamérica)
	}

	return &AzureTTSService{
		subscriptionKey: subscriptionKey,
		region:          region,
		logger:          logger,
		httpClient:      &http.Client{},
		voiceName:       voiceName,
		language:        language,
	}
}

// TextToSpeech convierte texto a audio usando Azure Speech Service
// Retorna los bytes del audio en formato MP3
// Soporta pausas usando [Pausa] o [Pausa: X] donde X es el número de segundos
func (s *AzureTTSService) TextToSpeech(text string) ([]byte, error) {
	if s.subscriptionKey == "" {
		return nil, fmt.Errorf("AZURE_SPEECH_KEY no configurada")
	}

	if s.region == "" {
		return nil, fmt.Errorf("AZURE_SPEECH_REGION no configurada")
	}

	if text == "" {
		return nil, fmt.Errorf("texto vacío")
	}

	// Preprocesar texto para mejorar pronunciación en español mexicano/latinoamericano
	text = s.preprocessForLatinAmericanSpanish(text)

	s.logger.Info(fmt.Sprintf("🎤 [Azure TTS] Generando audio con voz %s, idioma %s para texto: %s", s.voiceName, s.language, text))

	// Construir SSML con pausas procesadas
	// Azure TTS soporta etiquetas SSML <break> para pausas precisas
	ssmlContent := s.processPausesToSSML(text)
	ssml := fmt.Sprintf(`<speak version="1.0" xml:lang="%s"><voice name="%s">%s</voice></speak>`,
		s.language, s.voiceName, ssmlContent)
	ssmlBytes := []byte(ssml)

	// Endpoint de Azure Speech Service
	url := fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/v1", s.region)

	// Crear request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(ssmlBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Headers requeridos por Azure Speech Service
	req.Header.Set("Ocp-Apim-Subscription-Key", s.subscriptionKey)
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("X-Microsoft-OutputFormat", "audio-24khz-48kbitrate-mono-mp3") // Formato MP3 para WhatsApp

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
		s.logger.Error(fmt.Sprintf("Azure TTS API error (status %d): %s", resp.StatusCode, string(audioData)))
		return nil, fmt.Errorf("azure TTS API returned status %d: %s", resp.StatusCode, string(audioData))
	}

	s.logger.Info(fmt.Sprintf("✅ [Azure TTS] Audio generado: %d bytes", len(audioData)))
	return audioData, nil
}

// processPausesToSSML procesa las etiquetas de pausa y las convierte en etiquetas SSML <break>
// Azure TTS soporta etiquetas SSML <break> para pausas precisas
// Procesa tanto [Pausa] o [Pausa: X] como la palabra "pausa" como texto
func (s *AzureTTSService) processPausesToSSML(text string) string {
	// IMPORTANTE: Procesar pausas ANTES de escapar XML para que los regex funcionen correctamente

	// 1. Primero procesar la palabra "pausa" como texto (palabra completa, no parte de otra palabra)
	// Esto captura casos donde el texto dice "pausa" sin corchetes
	reWord := regexp.MustCompile(`(?i)\bpausa\b`) // (?i) para case-insensitive
	text = reWord.ReplaceAllString(text, `__PAUSA_WORD__`)

	// 2. Procesar etiquetas [Pausa] o [Pausa: X] (con corchetes)
	reTag := regexp.MustCompile(`\[Pausa(?::\s*(\d+))?\]`)

	text = reTag.ReplaceAllStringFunc(text, func(match string) string {
		parts := reTag.FindStringSubmatch(match)

		if len(parts) > 1 && parts[1] != "" {
			// Pausa con duración específica en segundos
			seconds := parts[1]
			return fmt.Sprintf(`__PAUSA_TAG_%s__`, seconds)
		}

		// Pausa corta por defecto (0.5 segundos)
		return `__PAUSA_TAG_0.5__`
	})

	// 3. Ahora escapar el texto para evitar problemas con caracteres especiales XML
	escapedText := escapeXMLText(text)

	// 4. Reemplazar los marcadores con las etiquetas SSML <break>
	// Reemplazar pausas de palabra
	escapedText = strings.ReplaceAll(escapedText, escapeXMLText(`__PAUSA_WORD__`), `<break time="0.5s"/>`)

	// Reemplazar pausas de etiquetas
	rePlaceholder := regexp.MustCompile(`__PAUSA_TAG_([0-9.]+)__`)
	escapedText = rePlaceholder.ReplaceAllStringFunc(escapedText, func(match string) string {
		parts := rePlaceholder.FindStringSubmatch(match)
		if len(parts) > 1 {
			seconds := parts[1]
			return fmt.Sprintf(`<break time="%ss"/>`, seconds)
		}
		return `<break time="0.5s"/>`
	})

	return escapedText
}

// preprocessForLatinAmericanSpanish preprocesa el texto para mejorar la pronunciación en español latinoamericano
func (s *AzureTTSService) preprocessForLatinAmericanSpanish(text string) string {
	// Normalizar espacios y puntuación
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	text = regexp.MustCompile(`\s*,\s*`).ReplaceAllString(text, ", ")
	text = regexp.MustCompile(`\s*\.\s*`).ReplaceAllString(text, ". ")

	// Asegurar que las preguntas tengan la entonación correcta
	text = regexp.MustCompile(`\?+`).ReplaceAllString(text, "?")
	text = regexp.MustCompile(`\!+`).ReplaceAllString(text, "!")

	// Para español latinoamericano, mantener las palabras comunes de la región
	// No necesitamos reemplazar palabras como en castellano

	return text
}

// escapeXMLText escapa caracteres especiales XML en el texto
func escapeXMLText(text string) string {
	// Escapar caracteres XML especiales
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	text = strings.ReplaceAll(text, "'", "&apos;")
	return text
}
