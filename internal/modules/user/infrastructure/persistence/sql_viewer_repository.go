package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/dukk308/golang-clean-arc/internal/modules/user/domain"
)

type ViewerRepository struct {
	db *gorm.DB
}

func NewViewerRepository(db *gorm.DB) domain.IViewerRepository {
	return &ViewerRepository{
		db: db,
	}
}

func (r *ViewerRepository) Create(ctx context.Context, viewer *domain.Viewer) error {
	return nil
}

func (r *ViewerRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *ViewerRepository) GetAll(ctx context.Context) ([]*domain.Viewer, error) {
	return nil, nil
}

func (r *ViewerRepository) GetByID(ctx context.Context, id string) (*domain.Viewer, error) {
	return nil, nil
}

func (r *ViewerRepository) Update(ctx context.Context, viewer *domain.Viewer) error {
	return nil
}
