package ports

import (
	"context"
	"whatsapp-api-go/internal/domain/entities"
)

// FlowRepository define las operaciones para flujos
type FlowRepository interface {
	// Save guarda un flujo
	Save(ctx context.Context, flow *entities.Flow) error

	// FindByID busca un flujo por ID
	FindByID(ctx context.Context, flowID string) (*entities.Flow, error)

	// FindDefault busca el flujo por defecto de una instancia
	FindDefault(ctx context.Context, instanceID string) (*entities.Flow, error)

	// FindByTenant busca flujos de un tenant
	FindByTenant(ctx context.Context, tenantID string) ([]*entities.Flow, error)

	// Update actualiza un flujo
	Update(ctx context.Context, flow *entities.Flow) error

	// Delete elimina un flujo
	Delete(ctx context.Context, flowID string) error
}

// FlowSessionRepository define las operaciones para sesiones de flujo
type FlowSessionRepository interface {
	// Save guarda una sesión
	Save(ctx context.Context, session *entities.FlowSession) error

	// FindByID busca una sesión por ID
	FindByID(ctx context.Context, sessionID string) (*entities.FlowSession, error)

	// FindActiveByConversation busca la sesión activa de una conversación
	FindActiveByConversation(ctx context.Context, conversationID string) (*entities.FlowSession, error)

	// Update actualiza una sesión
	Update(ctx context.Context, session *entities.FlowSession) error

	// FindInactiveSessions busca sesiones inactivas (para timeout)
	FindInactiveSessions(ctx context.Context, minutesInactive int) ([]*entities.FlowSession, error)
}

// FlowEngine define el motor de procesamiento de flujos
type FlowEngine interface {
	// StartFlow inicia un nuevo flujo para una conversación
	StartFlow(ctx context.Context, conversationID string, flowID string, tenantID string, instanceID string) (*entities.FlowSession, error)

	// ProcessMessage procesa un mensaje en el contexto de un flujo
	ProcessMessage(ctx context.Context, session *entities.FlowSession, message *entities.Message) error

	// ProcessNode procesa un nodo específico
	ProcessNode(ctx context.Context, session *entities.FlowSession, node *entities.FlowNode) error

	// MoveToNextNode avanza al siguiente nodo según los edges
	MoveToNextNode(ctx context.Context, session *entities.FlowSession, condition string) error
}


