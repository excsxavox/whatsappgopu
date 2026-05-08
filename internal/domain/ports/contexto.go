package ports

import (
	"context"
	"whatsapp-api-go/internal/domain/entities"
)

// ContextoRepository define el puerto para persistencia de contextos
type ContextoRepository interface {
	Save(ctx context.Context, contexto *entities.Contexto) error
	FindByID(ctx context.Context, contextoID string) (*entities.Contexto, error)
	FindByCompanyID(ctx context.Context, companyID string) (*entities.Contexto, error)
	FindActiveByCompanyID(ctx context.Context, companyID string) (*entities.Contexto, error)
	FindAll(ctx context.Context) ([]*entities.Contexto, error)
	Delete(ctx context.Context, contextoID string) error
}
