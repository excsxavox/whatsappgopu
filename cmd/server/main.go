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

	// 5.1 Inicializar Motor de Flujos
	log.Info("🔄 Configurando motor de flujos...")
	
	flowEngine := flow.NewFlowEngine(
		flowRepo,
		flowSessionRepo,
		whatsappAdapter,
		log,
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
