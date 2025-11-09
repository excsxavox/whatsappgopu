package ports

import (
	"context"
	"whatsapp-api-go/internal/domain/entities"
)

// SessionRepository define el puerto para persistencia de sesiones
// Este es un puerto de SALIDA (driving port)
type SessionRepository interface {
	// Save guarda o actualiza una sesión
	Save(ctx context.Context, session *entities.Session) error
	
	// FindByID busca una sesión por ID
	FindByID(ctx context.Context, sessionID string) (*entities.Session, error)
	
	// FindActive busca la sesión activa
	FindActive(ctx context.Context) (*entities.Session, error)
	
	// Delete elimina una sesión
	Delete(ctx context.Context, sessionID string) error
	
	// MarkAsInactive marca una sesión como inactiva
	MarkAsInactive(ctx context.Context, sessionID string) error
}

// SessionManager define el puerto para gestión de sesiones de WhatsApp
// Este es un puerto de SALIDA (driving port)
type SessionManager interface {
	// CreateSession crea una nueva sesión de WhatsApp
	CreateSession(ctx context.Context) (*entities.Session, error)
	
	// GetQRCode obtiene el código QR para autenticación
	GetQRCode(ctx context.Context) (string, error)
	
	// WaitForConnection espera a que se establezca la conexión
	WaitForConnection(ctx context.Context) (*entities.Session, error)
	
	// RestoreSession restaura una sesión existente
	RestoreSession(ctx context.Context, sessionID string) (*entities.Session, error)
	
	// DestroySession destruye una sesión
	DestroySession(ctx context.Context, sessionID string) error
}

