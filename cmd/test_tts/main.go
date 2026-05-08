package main

import (
	"context"
	"fmt"
	"os"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/infrastructure/adapters/whatsapp"
	"whatsapp-api-go/internal/infrastructure/services"
	"whatsapp-api-go/pkg/logger"
)

func main() {
	// Configuración desde argumentos o variables de entorno
	accessToken := os.Getenv("WABA_TOKEN")
	phoneNumberID := os.Getenv("WABA_PHONE_ID")
	apiVersion := os.Getenv("WABA_API_VERSION")
	toPhone := os.Getenv("TO_PHONE") // Número de destino (ej: 593992686734)
	googleTTSKey := os.Getenv("GOOGLE_TTS_API_KEY")
	openAIKey := os.Getenv("OPENAI_API_KEY")

	// Si no están en variables de entorno, usar argumentos
	if len(os.Args) > 1 {
		accessToken = os.Args[1]
	}
	if len(os.Args) > 2 {
		phoneNumberID = os.Args[2]
	}
	if len(os.Args) > 3 {
		toPhone = os.Args[3]
	}
	if len(os.Args) > 4 {
		apiVersion = os.Args[4]
	}

	// Validar parámetros requeridos
	if accessToken == "" {
		fmt.Println("❌ Error: Se requiere WABA_TOKEN o pasarlo como primer argumento")
		fmt.Println("Uso: go run cmd/test_tts/main.go <token> <phoneNumberID> <toPhone> [apiVersion]")
		fmt.Println("O usar variables de entorno: WABA_TOKEN, WABA_PHONE_ID, TO_PHONE")
		os.Exit(1)
	}

	if phoneNumberID == "" {
		fmt.Println("❌ Error: Se requiere WABA_PHONE_ID o pasarlo como segundo argumento")
		os.Exit(1)
	}

	if toPhone == "" {
		fmt.Println("❌ Error: Se requiere TO_PHONE o pasarlo como tercer argumento")
		os.Exit(1)
	}

	if apiVersion == "" {
		apiVersion = "v21.0" // Versión por defecto
	}

	fmt.Println("============================================")
	fmt.Println("🧪 Prueba Local: Google TTS + WhatsApp")
	fmt.Println("============================================")
	fmt.Println()

	// Inicializar logger
	log := logger.NewColorLogger()

	// 1. Inicializar Google TTS
	log.Info("🎤 Inicializando Google TTS...")
	googleTTSVoice := os.Getenv("GOOGLE_TTS_VOICE")
	googleTTSModel := os.Getenv("GOOGLE_TTS_MODEL")
	googleTTSLanguage := os.Getenv("GOOGLE_TTS_LANGUAGE")
	googleTTSPrompt := os.Getenv("GOOGLE_TTS_PROMPT")

	if googleTTSKey == "" {
		log.Warn("⚠️ GOOGLE_TTS_API_KEY no configurada, usando OpenAI TTS")
	}

	ttsService := services.NewUnifiedTTSService(
		googleTTSKey,
		openAIKey,
		log,
		googleTTSVoice,
		googleTTSModel,
		googleTTSLanguage,
		googleTTSPrompt,
	)

	// 2. Generar audio de prueba
	textoPrueba := "Hola, esta es una prueba de Google Text to Speech con acento castellano. ¿Puedes escucharme correctamente?"
	log.Info(fmt.Sprintf("📝 Generando audio para el texto: %s", textoPrueba))

	audioBytes, err := ttsService.TextToSpeech(textoPrueba)
	if err != nil {
		log.Error("❌ Error generando audio", "error", err)
		os.Exit(1)
	}

	log.Info(fmt.Sprintf("✅ Audio generado exitosamente (%d bytes)", len(audioBytes)))

	// 3. Inicializar adaptador de WhatsApp
	log.Info("📱 Inicializando adaptador de WhatsApp...")
	whatsappAdapter := whatsapp.NewCloudAPIAdapter(phoneNumberID, accessToken, apiVersion, log)

	// 4. Crear mensaje con audio
	ctx := context.Background()
	log.Info(fmt.Sprintf("📤 Enviando audio a WhatsApp (número: %s)...", toPhone))

	msg := &entities.Message{
		To:        toPhone,
		Direction: "out",
		MessageData: entities.MessageData{
			Type: "audio",
			Media: &entities.MediaContent{
				MimeType: "audio/mpeg",
				Data:     audioBytes,
			},
		},
	}

	err = whatsappAdapter.SendMessage(ctx, msg)
	if err != nil {
		log.Error("❌ Error enviando audio a WhatsApp", "error", err)
		os.Exit(1)
	}

	log.Info("✅ Audio enviado exitosamente a WhatsApp!")
	fmt.Println()
	fmt.Println("🎉 Prueba completada exitosamente!")
}
