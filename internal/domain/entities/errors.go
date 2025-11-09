package entities

import "errors"

var (
	// Errores de mensaje
	ErrInvalidRecipient = errors.New("recipient is invalid or empty")
	ErrEmptyMessage     = errors.New("message content is empty")
	ErrMessageNotFound  = errors.New("message not found")

	// Errores de sesión
	ErrSessionNotFound  = errors.New("session not found")
	ErrSessionNotActive = errors.New("session is not active")
	ErrSessionExpired   = errors.New("session has expired")
	ErrNotAuthenticated = errors.New("not authenticated")

	// Errores de conexión
	ErrNotConnected     = errors.New("not connected to WhatsApp")
	ErrConnectionFailed = errors.New("connection failed")
	ErrQRCodeExpired    = errors.New("QR code expired")
)
