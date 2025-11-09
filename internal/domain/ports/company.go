package ports

import (
	"context"
	"whatsapp-api-go/internal/domain/entities"
)

// CompanyRepository define el puerto para persistencia de empresas
type CompanyRepository interface {
	Save(ctx context.Context, company *entities.Company) error
	FindByID(ctx context.Context, companyID string) (*entities.Company, error)
	FindByCode(ctx context.Context, code string) (*entities.Company, error)
	FindAll(ctx context.Context) ([]*entities.Company, error)
	FindActive(ctx context.Context) ([]*entities.Company, error)
	FindInactive(ctx context.Context) ([]*entities.Company, error)
	Delete(ctx context.Context, companyID string) error
}
