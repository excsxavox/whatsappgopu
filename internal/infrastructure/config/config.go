package config

import "os"

// Config contiene la configuración para WhatsApp Business Cloud API
type Config struct {
	// API Server
	APIPort string

	// WhatsApp Cloud API
	VerifyToken   string // Token para validar webhook (hub.verify_token)
	AppSecret     string // App Secret para validar firma HMAC
	PhoneNumberID string // ID del número de teléfono (WABA_PHONE_ID)
	AccessToken   string // Token de acceso permanente (WABA_TOKEN)
	APIVersion    string // Versión de la API (v20.0)

	// MongoDB
	MongoURI string // URI de conexión a MongoDB
	MongoDB  string // Nombre de la base de datos

	// Rate Limiting
	PairRateLimit int // Mensajes por usuario cada N segundos (default: 1 msg / 6s)

	// Logs
	LogLevel string
}

// Load carga la configuración desde variables de entorno
func Load() *Config {
	// Permitir MONGO_URI o MONGODB_URL
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGODB_URL")
	}
	if mongoURI == "" {
		panic("Variable de entorno requerida: MONGO_URI o MONGODB_URL")
	}

	return &Config{
		// API Server
		APIPort: getEnv("API_PORT", "8080"),

		// WhatsApp Cloud API - REQUERIDOS
		VerifyToken:   getEnvRequired("WHATSAPP_VERIFY_TOKEN"),
		AppSecret:     getEnvRequired("WHATSAPP_APP_SECRET"),
		PhoneNumberID: getEnvRequired("WABA_PHONE_ID"),
		AccessToken:   getEnvRequired("WABA_TOKEN"),
		APIVersion:    getEnv("WABA_API_VERSION", "v20.0"),

		// MongoDB - REQUERIDOS
		MongoURI: mongoURI,
		MongoDB:  getEnv("MONGO_DB", "whatsapp_api"),

		// Rate Limiting
		PairRateLimit: 6, // 1 mensaje cada 6 segundos por usuario

		// Logs
		LogLevel: getEnv("LOG_LEVEL", "INFO"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Variable de entorno requerida: " + key)
	}
	return value
}
