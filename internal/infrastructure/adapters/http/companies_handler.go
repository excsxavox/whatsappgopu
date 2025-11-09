package http

import (
	"encoding/json"
	"net/http"
	"time"

	"whatsapp-api-go/internal/domain/entities"
	"whatsapp-api-go/internal/domain/ports"

	"github.com/google/uuid"
)

// CompaniesHandler maneja los endpoints de empresas
type CompaniesHandler struct {
	companyRepo ports.CompanyRepository
	logger      ports.Logger
}

// NewCompaniesHandler crea un nuevo handler de empresas
func NewCompaniesHandler(companyRepo ports.CompanyRepository, logger ports.Logger) *CompaniesHandler {
	return &CompaniesHandler{
		companyRepo: companyRepo,
		logger:      logger,
	}
}

// ListCompanies GET /api/companies
func (h *CompaniesHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	var companies []*entities.Company
	var err error

	switch status {
	case "active":
		companies, err = h.companyRepo.FindActive(r.Context())
	case "inactive":
		companies, err = h.companyRepo.FindInactive(r.Context())
	default:
		companies, err = h.companyRepo.FindAll(r.Context())
	}

	if err != nil {
		h.logger.Error("Error al listar empresas", "error", err)
		http.Error(w, "Error al listar empresas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"companies": companies,
		"count":     len(companies),
	})
}

// GetCompany GET /api/companies/{id}
func (h *CompaniesHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	// Extraer ID de la URL (implementación simple, mejorar con router)
	id := r.URL.Path[len("/api/companies/"):]

	company, err := h.companyRepo.FindByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Error al buscar empresa", "error", err)
		http.Error(w, "Error al buscar empresa", http.StatusInternalServerError)
		return
	}

	if company == nil {
		http.Error(w, "Empresa no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

// CreateCompany POST /api/companies
func (h *CompaniesHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code           string `json:"code"`
		Name           string `json:"name"`
		BusinessType   string `json:"business_type"`
		WhatsAppNumber string `json:"whatsapp_number"`
		PhoneNumberID  string `json:"phone_number_id"`
		AccessToken    string `json:"access_token"`
		WebhookToken   string `json:"webhook_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Verificar si ya existe
	existing, _ := h.companyRepo.FindByCode(r.Context(), req.Code)
	if existing != nil {
		http.Error(w, "Ya existe una empresa con ese código", http.StatusConflict)
		return
	}

	// Crear empresa
	company := entities.NewCompany(req.Code, req.Name, req.BusinessType, req.WhatsAppNumber)
	company.ID = uuid.New().String()

	if req.PhoneNumberID != "" && req.AccessToken != "" {
		company.UpdateMetaCredentials(req.PhoneNumberID, req.AccessToken, req.WebhookToken)
		company.Activate() // Auto-activar si tiene credenciales
	}

	if err := h.companyRepo.Save(r.Context(), company); err != nil {
		h.logger.Error("Error al crear empresa", "error", err)
		http.Error(w, "Error al crear empresa", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Empresa creada", "id", company.ID, "code", company.Code)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"company": company,
		"message": "Empresa creada exitosamente",
	})
}

// UpdateCompany PUT /api/companies/{id}
func (h *CompaniesHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/companies/"):]

	var req struct {
		Name           string `json:"name"`
		BusinessType   string `json:"business_type"`
		WhatsAppNumber string `json:"whatsapp_number"`
		PhoneNumberID  string `json:"phone_number_id"`
		AccessToken    string `json:"access_token"`
		WebhookToken   string `json:"webhook_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	company, err := h.companyRepo.FindByID(r.Context(), id)
	if err != nil || company == nil {
		http.Error(w, "Empresa no encontrada", http.StatusNotFound)
		return
	}

	// Actualizar campos
	if req.Name != "" {
		company.Name = req.Name
	}
	if req.BusinessType != "" {
		company.BusinessType = req.BusinessType
	}
	if req.WhatsAppNumber != "" {
		company.WhatsAppNumber = req.WhatsAppNumber
	}
	if req.PhoneNumberID != "" {
		company.UpdateMetaCredentials(req.PhoneNumberID, req.AccessToken, req.WebhookToken)
	}

	company.UpdatedAt = time.Now()

	if err := h.companyRepo.Save(r.Context(), company); err != nil {
		h.logger.Error("Error al actualizar empresa", "error", err)
		http.Error(w, "Error al actualizar empresa", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"company": company,
		"message": "Empresa actualizada exitosamente",
	})
}

// ActivateCompany POST /api/companies/{id}/activate
func (h *CompaniesHandler) ActivateCompany(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/companies/"):]
	id = id[:len(id)-len("/activate")]

	company, err := h.companyRepo.FindByID(r.Context(), id)
	if err != nil || company == nil {
		http.Error(w, "Empresa no encontrada", http.StatusNotFound)
		return
	}

	company.Activate()

	if err := h.companyRepo.Save(r.Context(), company); err != nil {
		h.logger.Error("Error al activar empresa", "error", err)
		http.Error(w, "Error al activar empresa", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Empresa activada", "id", company.ID, "code", company.Code)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"company": company,
		"message": "Empresa activada exitosamente",
	})
}

// DeactivateCompany POST /api/companies/{id}/deactivate
func (h *CompaniesHandler) DeactivateCompany(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/companies/"):]
	id = id[:len(id)-len("/deactivate")]

	company, err := h.companyRepo.FindByID(r.Context(), id)
	if err != nil || company == nil {
		http.Error(w, "Empresa no encontrada", http.StatusNotFound)
		return
	}

	company.Deactivate()

	if err := h.companyRepo.Save(r.Context(), company); err != nil {
		h.logger.Error("Error al desactivar empresa", "error", err)
		http.Error(w, "Error al desactivar empresa", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Empresa desactivada", "id", company.ID, "code", company.Code)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"company": company,
		"message": "Empresa desactivada exitosamente",
	})
}

// DeleteCompany DELETE /api/companies/{id}
func (h *CompaniesHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/companies/"):]

	if err := h.companyRepo.Delete(r.Context(), id); err != nil {
		h.logger.Error("Error al eliminar empresa", "error", err)
		http.Error(w, "Error al eliminar empresa", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Empresa eliminada", "id", id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Empresa eliminada exitosamente",
	})
}
