package http

import (
	"fmt"
	"net/http"

	"whatsapp-api-go/internal/domain/ports"
)

// Server representa el servidor HTTP
type Server struct {
	httpAdapter      *HTTPAdapter
	webhookHandler   *WebhookHandler
	companiesHandler *CompaniesHandler
	port             string
	logger           ports.Logger
}

// NewServer crea un nuevo servidor HTTP
func NewServer(
	httpAdapter *HTTPAdapter,
	webhookHandler *WebhookHandler,
	companiesHandler *CompaniesHandler,
	port string,
	logger ports.Logger,
) *Server {
	return &Server{
		httpAdapter:      httpAdapter,
		webhookHandler:   webhookHandler,
		companiesHandler: companiesHandler,
		port:             port,
		logger:           logger,
	}
}

// Start inicia el servidor HTTP
func (s *Server) Start() error {
	// CORS middleware simple
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Configurar CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Enrutamiento manual
		s.route(w, r)
	})

	s.logger.Info(fmt.Sprintf("🚀 Servidor API iniciado en http://localhost:%s", s.port))
	s.logger.Info("📡 Webhook disponible en: GET/POST /webhook")
	s.logger.Info("🏢 API Empresas disponible en: /api/companies")
	return http.ListenAndServe(":"+s.port, nil)
}

// route maneja el enrutamiento de las peticiones
func (s *Server) route(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method

	// Endpoints básicos
	if path == "/health" && method == "GET" {
		s.httpAdapter.HealthHandler(w, r)
		return
	}
	if path == "/status" && method == "GET" {
		s.httpAdapter.StatusHandler(w, r)
		return
	}
	if path == "/send" && method == "POST" {
		s.httpAdapter.SendMessageHandler(w, r)
		return
	}

	// Webhooks de Meta
	if path == "/webhook" {
		if method == "GET" {
			s.webhookHandler.VerifyWebhook(w, r)
		} else if method == "POST" {
			s.webhookHandler.ReceiveWebhook(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// API de Empresas
	if path == "/api/companies" && method == "GET" {
		s.companiesHandler.ListCompanies(w, r)
		return
	}
	if path == "/api/companies" && method == "POST" {
		s.companiesHandler.CreateCompany(w, r)
		return
	}
	if len(path) > 15 && path[:15] == "/api/companies/" {
		// /api/companies/{id}
		if method == "GET" {
			s.companiesHandler.GetCompany(w, r)
			return
		}
		if method == "PUT" {
			s.companiesHandler.UpdateCompany(w, r)
			return
		}
		if method == "DELETE" {
			s.companiesHandler.DeleteCompany(w, r)
			return
		}
		// /api/companies/{id}/activate
		if len(path) > 24 && path[len(path)-9:] == "/activate" && method == "POST" {
			s.companiesHandler.ActivateCompany(w, r)
			return
		}
		// /api/companies/{id}/deactivate
		if len(path) > 26 && path[len(path)-11:] == "/deactivate" && method == "POST" {
			s.companiesHandler.DeactivateCompany(w, r)
			return
		}
	}

	// 404
	http.NotFound(w, r)
}
