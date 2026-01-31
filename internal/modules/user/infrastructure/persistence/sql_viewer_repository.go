package persistence

import (
	"context"

	"gorm.io/gorm"

	"github.com/dukk308/golang-clean-arch-starter/internal/modules/user/domain"
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
	sqlUser := &SQLUser{}
	sqlUser.FromDomainViewer(viewer)
	return r.db.WithContext(ctx).Create(sqlUser).Error
}

func (r *ViewerRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&SQLUser{}).Error
}

func (r *ViewerRepository) GetAll(ctx context.Context) ([]*domain.Viewer, error) {
	var sqlUsers []SQLUser
	if err := r.db.WithContext(ctx).Find(&sqlUsers).Error; err != nil {
		return nil, err
	}

	viewers := make([]*domain.Viewer, len(sqlUsers))
	for i, sqlUser := range sqlUsers {
		viewers[i] = sqlUser.ToDomainViewer()
	}

	return viewers, nil
}

func (r *ViewerRepository) GetByEmail(ctx context.Context, email string) (*domain.Viewer, error) {
	var sqlUser SQLUser
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&sqlUser).Error; err != nil {
		return nil, err
	}

	viewer := sqlUser.ToDomainViewer()
	return viewer, nil
}

func (r *ViewerRepository) GetByID(ctx context.Context, id string) (*domain.Viewer, error) {
	var sqlUser SQLUser
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&sqlUser).Error; err != nil {
		return nil, err
	}

	viewer := sqlUser.ToDomainViewer()
	return viewer, nil
}

func (r *ViewerRepository) Update(ctx context.Context, viewer *domain.Viewer) error {
	sqlUser := &SQLUser{}
	sqlUser.FromDomainViewer(viewer)
	return r.db.WithContext(ctx).Save(sqlUser).Error
}
