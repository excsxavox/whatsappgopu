package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	// Application layer
	"whatsapp-api-go/internal/application/usecases"

	// Infrastructure layer
	"whatsapp-api-go/internal/infrastructure/adapters/http"
	"whatsapp-api-go/internal/infrastructure/adapters/storage"
	"whatsapp-api-go/internal/infrastructure/adapters/whatsapp"
	"whatsapp-api-go/internal/infrastructure/config"
	"whatsapp-api-go/internal/infrastructure/flow"
	"whatsapp-api-go/internal/infrastructure/services"

	// Utilities
	"whatsapp-api-go/pkg/logger"
)

func main() {
	fmt.Println("============================================")
	fmt.Println("🚀 WhatsApp Business Cloud API - Server")
	fmt.Println("   Arquitectura Hexagonal + MongoDB")
	fmt.Println("============================================")
	fmt.Println()

	// 1. Cargar configuración
	fmt.Println("📋 Cargando configuración...")
	cfg := config.Load()
	fmt.Printf("   ✓ API Port: %s\n", cfg.APIPort)
	fmt.Printf("   ✓ Phone Number ID: %s\n", cfg.PhoneNumberID)
	fmt.Printf("   ✓ API Version: %s\n", cfg.APIVersion)
	fmt.Printf("   ✓ MongoDB: %s\n", cfg.MongoDB)
	fmt.Println()

	// 2. Inicializar Logger
	log := logger.NewColorLogger()

	// 3. Conectar a MongoDB
	log.Info("📊 Conectando a MongoDB...")
	ctx := context.Background()

	mongoClient, err := storage.NewMongoClient(ctx, cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Error("❌ Error conectando a MongoDB", "error", err)
		os.Exit(1)
	}
	defer mongoClient.Close(ctx)

	log.Info("✅ Conectado a MongoDB")

	// Crear índices
	log.Info("📑 Creando índices...")
	if err := mongoClient.CreateIndexes(ctx); err != nil {
		log.Warn("⚠️  Error creando índices", "error", err)
	} else {
		log.Info("✅ Índices creados")
	}
	fmt.Println()

	// 4. Inicializar Adaptadores de Infraestructura
	log.Info("📦 Inicializando adaptadores...")

	// Adaptador de WhatsApp Cloud API
	whatsappAdapter := whatsapp.NewCloudAPIAdapter(
		cfg.PhoneNumberID,
		cfg.AccessToken,
		cfg.APIVersion,
		log,
	)

	// Repositorios MongoDB
	messageRepo := storage.NewMongoMessageRepository(mongoClient.GetDatabase())
	sessionRepo := storage.NewMongoSessionRepository(mongoClient.GetDatabase())
	companyRepo := storage.NewMongoCompanyRepository(mongoClient.GetDatabase())
	flowRepo := storage.NewMongoFlowRepository(mongoClient)
	flowSessionRepo := storage.NewMongoFlowSessionRepository(mongoClient)
	contextoRepo := storage.NewMongoContextoRepository(mongoClient.GetDatabase())

	log.Info("✅ Adaptadores inicializados")
	fmt.Println()

	// 5. Inicializar Casos de Uso (Application Layer)
	log.Info("🔧 Configurando casos de uso...")

	sendMessageUseCase := usecases.NewSendMessageUseCase(
		whatsappAdapter,
		messageRepo,
		log,
	)

	getConnectionStatusUseCase := usecases.NewGetConnectionStatusUseCase(
		whatsappAdapter,
		log,
	)

	establishConnectionUseCase := usecases.NewEstablishConnectionUseCase(
		whatsappAdapter,
		sessionRepo,
		log,
	)

	disconnectUseCase := usecases.NewDisconnectUseCase(
		whatsappAdapter,
		sessionRepo,
		log,
	)

	handleWebhookUseCase := usecases.NewHandleWebhookUseCase(
		sendMessageUseCase,
		log,
	)

	log.Info("✅ Casos de uso configurados")
	fmt.Println()

	// 5.1 Inicializar Servicio de Transcripción
	log.Info("🎤 Configurando servicio de transcripción...")
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		log.Warn("⚠️ OPENAI_API_KEY no configurada, transcripción de audio deshabilitada")
	}
	transcriptionService := services.NewTranscriptionService(openAIKey, log)

	// 5.1.1 Inicializar Servicio de Validación con IA
	log.Info("🤖 Configurando servicio de validación con IA...")
	aiValidationService := services.NewAIValidationService(openAIKey, log)

	// 5.1.1.1 Inicializar Servicio de IA con Contexto
	log.Info("🧠 Configurando servicio de IA con contexto...")
	aiContextService := services.NewAIContextService(openAIKey, log)

	// 5.1.2 Inicializar Servicio de Texto a Voz (TTS)
	log.Info("🔊 Configurando servicio de texto a voz (TTS)...")
	// Prioridad: Azure > Google > OpenAI
	// Azure Speech Service (máxima prioridad)
	azureSpeechKey := os.Getenv("AZURE_SPEECH_KEY")       // Clave de suscripción de Azure Speech Service
	azureSpeechRegion := os.Getenv("AZURE_SPEECH_REGION")  // Región de Azure (ej: "eastus", "westus")
	azureTTSVoice := os.Getenv("AZURE_TTS_VOICE")         // Opcional: ej. "es-ES-ElviraNeural" (default: "es-ES-ElviraNeural")
	azureTTSLanguage := os.Getenv("AZURE_TTS_LANGUAGE")    // Opcional: "es-ES" o "es-MX" (default: "es-ES")

	// Google TTS (segunda prioridad)
	googleTTSKey := os.Getenv("GOOGLE_TTS_API_KEY")               // Clave de API (para modelo standard)
	googleTTSJSON := os.Getenv("GOOGLE_TTS_SERVICE_ACCOUNT_JSON") // JSON de cuenta de servicio (para Gemini Pro TTS)
	googleTTSVoice := os.Getenv("GOOGLE_TTS_VOICE")               // Opcional: ej. "es-ES-Neural2-A" para modelo standard
	googleTTSModel := os.Getenv("GOOGLE_TTS_MODEL")               // Opcional: "gemini" o "standard" (default: "standard")
	googleTTSLanguage := os.Getenv("GOOGLE_TTS_LANGUAGE")         // Opcional: "es-ES" o "es-MX" (default: "es-ES")
	googleTTSPrompt := os.Getenv("GOOGLE_TTS_PROMPT")             // Opcional: prompt para controlar el tono (solo para Gemini)

	// Determinar qué autenticación usar según el modelo de Google
	// Si el modelo es "standard", usar API key (más simple, no requiere Vertex AI)
	// Si el modelo es "gemini", usar JSON de cuenta de servicio (requiere Vertex AI)
	var googleTTSAuth string
	if googleTTSModel == "" || googleTTSModel == "standard" {
		// Modelo standard: priorizar API key
		googleTTSAuth = googleTTSKey
		if googleTTSAuth == "" {
			googleTTSAuth = googleTTSJSON // Fallback a JSON si no hay API key
		}
	} else {
		// Modelo gemini: requiere JSON de cuenta de servicio
		googleTTSAuth = googleTTSJSON
		if googleTTSAuth == "" {
			googleTTSAuth = googleTTSKey // Fallback a API key si no hay JSON
		}
	}

	// Determinar voz e idioma a usar (priorizar Azure, luego Google, luego defaults)
	voiceName := azureTTSVoice
	if voiceName == "" {
		voiceName = googleTTSVoice
	}
	language := azureTTSLanguage
	if language == "" {
		language = googleTTSLanguage
	}
	if language == "" {
		language = "es-ES" // Default: español de España
	}

	// Log de configuración para debugging
	if azureSpeechKey != "" && azureSpeechRegion != "" {
		log.Info(fmt.Sprintf("✅ AZURE_SPEECH_KEY encontrada (longitud: %d)", len(azureSpeechKey)))
		log.Info(fmt.Sprintf("✅ AZURE_SPEECH_REGION: %s", azureSpeechRegion))
		log.Info(fmt.Sprintf("   Voz: %s, Idioma: %s", voiceName, language))
	} else if googleTTSAuth != "" {
		if googleTTSJSON != "" {
			log.Info("✅ GOOGLE_TTS_SERVICE_ACCOUNT_JSON encontrada (autenticación con cuenta de servicio)")
		} else if googleTTSKey != "" {
			log.Info(fmt.Sprintf("✅ GOOGLE_TTS_API_KEY encontrada (longitud: %d)", len(googleTTSKey)))
		}
		log.Info(fmt.Sprintf("   Voz: %s, Modelo: %s, Idioma: %s", googleTTSVoice, googleTTSModel, googleTTSLanguage))
	} else {
		log.Info("ℹ️ AZURE_SPEECH_KEY o GOOGLE_TTS_API_KEY no configurada, usando OpenAI TTS")
	}
	ttsService := services.NewUnifiedTTSService(azureSpeechKey, azureSpeechRegion, googleTTSAuth, openAIKey, log, voiceName, googleTTSModel, language, googleTTSPrompt)

	// 5.2 Inicializar Motor de Flujos
	log.Info("🔄 Configurando motor de flujos...")

	flowEngine := flow.NewFlowEngine(
		flowRepo,
		flowSessionRepo,
		whatsappAdapter,
		log,
		transcriptionService,
		aiValidationService,
		aiContextService,
		contextoRepo,
		ttsService,
	)

	startFlowUseCase := usecases.NewStartFlowUseCase(
		flowEngine,
		flowRepo,
		log,
	)

	processFlowMessageUseCase := usecases.NewProcessFlowMessageUseCase(
		flowEngine,
		flowSessionRepo,
		log,
	)

	log.Info("✅ Motor de flujos configurado")
	fmt.Println()

	// 6. Verificar conexión con Cloud API
	log.Info("📱 Verificando Cloud API...")

	connection, err := establishConnectionUseCase.Execute(ctx)
	if err != nil {
		log.Error("❌ Error al verificar Cloud API", "error", err)
		os.Exit(1)
	}

	log.Info("✅ Cloud API verificada y lista")
	log.Info(fmt.Sprintf("   Estado: conectado=%v, autenticado=%v",
		connection.IsConnected, connection.IsLoggedIn))
	fmt.Println()

	// 7. Inicializar Adaptadores HTTP
	log.Info("🌐 Configurando servidor HTTP...")

	// Adaptador HTTP para endpoints REST
	httpAdapter := http.NewHTTPAdapter(
		sendMessageUseCase,
		getConnectionStatusUseCase,
		handleWebhookUseCase,
		log,
	)

	// Handler de Webhooks de Meta (con flujos integrados)
	webhookHandler := http.NewWebhookHandler(
		cfg.VerifyToken,
		cfg.AppSecret,
		cfg.PhoneNumberID, // instanceID
		sendMessageUseCase,
		messageRepo,
		log,
		startFlowUseCase,
		processFlowMessageUseCase,
	)

	// Handler de Empresas
	companiesHandler := http.NewCompaniesHandler(
		companyRepo,
		log,
	)

	// Servidor HTTP
	httpServer := http.NewServer(
		httpAdapter,
		webhookHandler,
		companiesHandler,
		cfg.APIPort,
		log,
	)

	// 8. Mostrar endpoints disponibles
	fmt.Printf("📡 Endpoints REST:\n")
	fmt.Printf("   - GET  http://localhost:%s/health\n", cfg.APIPort)
	fmt.Printf("   - GET  http://localhost:%s/status\n", cfg.APIPort)
	fmt.Printf("   - POST http://localhost:%s/send\n", cfg.APIPort)
	fmt.Println()
	fmt.Printf("🏢 API Empresas:\n")
	fmt.Printf("   - GET    http://localhost:%s/api/companies\n", cfg.APIPort)
	fmt.Printf("   - POST   http://localhost:%s/api/companies\n", cfg.APIPort)
	fmt.Printf("   - GET    http://localhost:%s/api/companies/{id}\n", cfg.APIPort)
	fmt.Printf("   - PUT    http://localhost:%s/api/companies/{id}\n", cfg.APIPort)
	fmt.Printf("   - DELETE http://localhost:%s/api/companies/{id}\n", cfg.APIPort)
	fmt.Printf("   - POST   http://localhost:%s/api/companies/{id}/activate\n", cfg.APIPort)
	fmt.Printf("   - POST   http://localhost:%s/api/companies/{id}/deactivate\n", cfg.APIPort)
	fmt.Println()
	fmt.Printf("🔔 Webhooks Meta:\n")
	fmt.Printf("   - GET  http://localhost:%s/webhook (verificación)\n", cfg.APIPort)
	fmt.Printf("   - POST http://localhost:%s/webhook (eventos)\n", cfg.APIPort)
	fmt.Println()

	log.Info("⚠️  IMPORTANTE: Configura este webhook en Meta:")
	log.Info(fmt.Sprintf("   URL: https://tu-dominio.com/webhook"))
	log.Info(fmt.Sprintf("   Verify Token: %s", cfg.VerifyToken))
	fmt.Println()

	// 9. Manejar señales del sistema para graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Info("\n\n🛑 Señal de interrupción recibida...")
		log.Info("🔌 Cerrando conexiones...")
		_ = disconnectUseCase.Execute(context.Background())
		_ = mongoClient.Close(context.Background())
		cancel()
	}()

	// 10. Iniciar servidor HTTP en goroutine
	go func() {
		if err := httpServer.Start(); err != nil {
			log.Error("❌ Error al iniciar el servidor HTTP", "error", err)
			cancel()
		}
	}()

	log.Info("✅ Sistema iniciado correctamente")
	log.Info("⌛ Servidor escuchando webhooks de Meta...")
	log.Info("💾 Persistencia: MongoDB")
	fmt.Println()

	// 11. Mantener la aplicación corriendo
	<-ctx.Done()

	log.Info("✅ Aplicación cerrada correctamente")
}
