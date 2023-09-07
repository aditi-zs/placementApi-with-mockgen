package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/aditi-zs/Placement-API/entities"
)

type StudentStore interface {
	GetWithCompany(ctx context.Context, name string, branch string) ([]entities.Student, error)
	Get(ctx context.Context, name string, branch string) ([]entities.Student, error)
	GetByID(ctx context.Context, id uuid.UUID) (entities.Student, error)
	Create(ctx context.Context, stu *entities.Student) (entities.Student, error)
	Update(ctx context.Context, id uuid.UUID, stu *entities.Student) (entities.Student, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetCompanyByID(ctx context.Context, id uuid.UUID) (entities.Company, error)
}
type CompanyStore interface {
	Get(ctx context.Context) ([]entities.Company, error)
	GetByID(ctx context.Context, id uuid.UUID) (entities.Company, error)
	Create(ctx context.Context, cmp entities.Company) (entities.Company, error)
	Update(ctx context.Context, id uuid.UUID, cmp entities.Company) (entities.Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
