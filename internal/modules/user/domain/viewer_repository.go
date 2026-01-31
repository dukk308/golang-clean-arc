package domain

import "context"

type IViewerRepository interface {
	GetByID(ctx context.Context, id string) (*Viewer, error)
	GetByEmail(ctx context.Context, email string) (*Viewer, error)
	GetAll(ctx context.Context) ([]*Viewer, error)
	Create(ctx context.Context, viewer *Viewer) error
	Update(ctx context.Context, viewer *Viewer) error
	Delete(ctx context.Context, id string) error
}
