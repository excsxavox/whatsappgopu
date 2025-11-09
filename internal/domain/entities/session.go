package entities

import "time"

// Session representa una sesión de WhatsApp
type Session struct {
	ID             string
	PhoneNumber    string
	DeviceName     string
	IsActive       bool
	IsConnected    bool
	ConnectedAt    *time.Time
	LastSeen       *time.Time
	DisconnectedAt *time.Time
}

// NewSession crea una nueva sesión
func NewSession(phoneNumber string) *Session {
	return &Session{
		PhoneNumber: phoneNumber,
		IsActive:    false,
		IsConnected: false,
	}
}

// Connect marca la sesión como conectada
func (s *Session) Connect() {
	now := time.Now()
	s.IsConnected = true
	s.IsActive = true
	s.ConnectedAt = &now
	s.LastSeen = &now
}

// Disconnect marca la sesión como desconectada
func (s *Session) Disconnect() {
	now := time.Now()
	s.IsConnected = false
	s.DisconnectedAt = &now
}

// UpdateLastSeen actualiza la última vez que se vio activa
func (s *Session) UpdateLastSeen() {
	now := time.Now()
	s.LastSeen = &now
}
