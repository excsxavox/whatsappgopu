package entities

import "time"

// Company representa una empresa/negocio en el dominio
type Company struct {
	ID             string     `json:"id" bson:"_id"`
	Code           string     `json:"code" bson:"code"`
	Name           string     `json:"name" bson:"name"`
	BusinessType   string     `json:"business_type" bson:"business_type"` // Tipo de organización
	WhatsAppNumber string     `json:"whatsapp_number" bson:"whatsapp_number"`
	PhoneNumberID  string     `json:"phone_number_id" bson:"phone_number_id"`
	AccessToken    string     `json:"-" bson:"access_token"`  // No exponer en JSON
	WebhookToken   string     `json:"-" bson:"webhook_token"` // No exponer en JSON
	IsActive       bool       `json:"is_active" bson:"is_active"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" bson:"updated_at"`
	InvalidDate    *time.Time `json:"invalid_date,omitempty" bson:"invalid_date,omitempty"`
}

// NewCompany crea una nueva empresa
func NewCompany(code, name, businessType, whatsappNumber string) *Company {
	now := time.Now()
	return &Company{
		Code:           code,
		Name:           name,
		BusinessType:   businessType,
		WhatsAppNumber: whatsappNumber,
		IsActive:       false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Activate activa la empresa
func (c *Company) Activate() {
	c.IsActive = true
	c.UpdatedAt = time.Now()
	c.InvalidDate = nil
}

// Deactivate desactiva la empresa
func (c *Company) Deactivate() {
	c.IsActive = false
	now := time.Now()
	c.UpdatedAt = now
	c.InvalidDate = &now
}

// UpdateMetaCredentials actualiza las credenciales de Meta
func (c *Company) UpdateMetaCredentials(phoneNumberID, accessToken, webhookToken string) {
	c.PhoneNumberID = phoneNumberID
	c.AccessToken = accessToken
	c.WebhookToken = webhookToken
	c.UpdatedAt = time.Now()
}
