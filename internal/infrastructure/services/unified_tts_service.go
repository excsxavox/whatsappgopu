package services

import (
	"fmt"
	"whatsapp-api-go/internal/domain/ports"
)

// UnifiedTTSService servicio unificado que puede usar Azure TTS, Google TTS o OpenAI TTS
type UnifiedTTSService struct {
	azureService   *AzureTTSService
	googleService  *GoogleTTSService
	openAIService  *TTSService
	logger         ports.Logger
	preferred      string // "azure", "google" o "openai"
}

// NewUnifiedTTSService crea un nuevo servicio TTS unificado
// Prioridad: Azure > Google > OpenAI
func NewUnifiedTTSService(azureKey string, azureRegion string, googleAPIKey string, openAIKey string, logger ports.Logger, voiceName string, model string, language string, prompt string) *UnifiedTTSService {
	var azureService *AzureTTSService
	var googleService *GoogleTTSService
	var openAIService *TTSService
	preferred := "openai" // Por defecto usa OpenAI

	// Log de diagnóstico
	logger.Info(fmt.Sprintf("🔍 [TTS Init] Azure Speech Key presente: %v (longitud: %d)", azureKey != "", len(azureKey)))
	logger.Info(fmt.Sprintf("🔍 [TTS Init] Azure Speech Region: %s", azureRegion))
	logger.Info(fmt.Sprintf("🔍 [TTS Init] Google API Key presente: %v (longitud: %d)", googleAPIKey != "", len(googleAPIKey)))
	logger.Info(fmt.Sprintf("🔍 [TTS Init] OpenAI API Key presente: %v", openAIKey != ""))
	logger.Info(fmt.Sprintf("🔍 [TTS Init] Configuración - Voz: %s, Modelo: %s, Idioma: %s", voiceName, model, language))

	// Si Azure TTS está configurado, usarlo EXCLUSIVAMENTE (sin fallback)
	if azureKey != "" && azureRegion != "" {
		azureService = NewAzureTTSService(azureKey, azureRegion, logger, voiceName, language)
		preferred = "azure"
		logger.Info(fmt.Sprintf("✅ [TTS Init] Azure TTS configurado como servicio EXCLUSIVO (voz: %s, idioma: %s, región: %s)", voiceName, language, azureRegion))
		logger.Info("ℹ️ [TTS Init] Solo se usará Azure TTS, sin fallback a otros servicios")
	} else {
		// Solo configurar Google y OpenAI si Azure NO está configurado
		// Si Google TTS está configurado, usarlo como preferido si Azure no está disponible
		if googleAPIKey != "" {
			googleService = NewGoogleTTSService(googleAPIKey, logger, voiceName, model, language, prompt)
			preferred = "google"
			logger.Info(fmt.Sprintf("✅ [TTS Init] Google TTS configurado como servicio preferido (modelo: %s, voz: %s, idioma: %s)", model, voiceName, language))
		}

		// Si OpenAI TTS está configurado, usarlo como fallback o principal
		if openAIKey != "" {
			openAIService = NewTTSService(openAIKey, logger)
			if preferred == "openai" {
				logger.Info("✅ [TTS Init] OpenAI TTS configurado como servicio principal")
			} else {
				logger.Info("✅ [TTS Init] OpenAI TTS configurado como servicio de respaldo")
			}
		}
	}

	if azureService == nil && googleService == nil && openAIService == nil {
		logger.Warn("⚠️ [TTS Init] Ningún servicio TTS configurado (ni Azure, ni Google ni OpenAI)")
	}

	return &UnifiedTTSService{
		azureService:  azureService,
		googleService: googleService,
		openAIService: openAIService,
		logger:        logger,
		preferred:     preferred,
	}
}

// TextToSpeech convierte texto a audio usando el servicio TTS disponible
// Retorna los bytes del audio en formato MP3
// Prioridad: Azure > Google > OpenAI
func (s *UnifiedTTSService) TextToSpeech(text string) ([]byte, error) {
	// Log del estado del servicio
	s.logger.Info(fmt.Sprintf("🔍 [TTS] Servicio preferido: %s, Azure disponible: %v, Google disponible: %v, OpenAI disponible: %v", 
		s.preferred, s.azureService != nil, s.googleService != nil, s.openAIService != nil))
	
	// Si Azure está configurado, usarlo EXCLUSIVAMENTE (sin fallback)
	if s.preferred == "azure" && s.azureService != nil {
		s.logger.Info("🎤 [TTS] Usando Azure TTS para generar audio (modo exclusivo)")
		audio, err := s.azureService.TextToSpeech(text)
		if err != nil {
			s.logger.Error(fmt.Sprintf("❌ [TTS] Error con Azure TTS: %v", err))
			return nil, fmt.Errorf("error con Azure TTS (servicio exclusivo): %w", err)
		}
		s.logger.Info(fmt.Sprintf("✅ [TTS] Audio generado exitosamente con Azure TTS (%d bytes)", len(audio)))
		return audio, nil
	}

	// Solo usar Google/OpenAI si Azure NO está configurado
	// Intentar usar Google TTS si está disponible
	if s.googleService != nil {
		s.logger.Info("🎤 [TTS] Intentando usar Google TTS para generar audio")
		audio, err := s.googleService.TextToSpeech(text)
		if err == nil {
			s.logger.Info(fmt.Sprintf("✅ [TTS] Audio generado exitosamente con Google TTS (%d bytes)", len(audio)))
			return audio, nil
		}
		s.logger.Error(fmt.Sprintf("❌ [TTS] Error con Google TTS: %v", err))
		s.logger.Warn("⚠️ [TTS] Intentando fallback con OpenAI TTS")
		// Si falla, intentar con OpenAI como fallback
	}

	// Usar OpenAI TTS si está disponible
	if s.openAIService != nil {
		s.logger.Info("🎤 [TTS] Usando OpenAI TTS para generar audio")
		audio, err := s.openAIService.TextToSpeech(text)
		if err != nil {
			s.logger.Error(fmt.Sprintf("❌ [TTS] Error con OpenAI TTS: %v", err))
			return nil, err
		}
		s.logger.Info(fmt.Sprintf("✅ [TTS] Audio generado exitosamente con OpenAI TTS (%d bytes)", len(audio)))
		return audio, nil
	}

	return nil, fmt.Errorf("ningún servicio TTS disponible")
}

