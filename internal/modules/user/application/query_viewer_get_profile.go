package application

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/user/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
)

type ViewerGetProfileQuery struct {
	repository domain.IViewerRepository
}

func NewViewerGetProfileQuery(repository domain.IViewerRepository) *ViewerGetProfileQuery {
	return &ViewerGetProfileQuery{
		repository: repository,
	}
}

func (q *ViewerGetProfileQuery) Execute(ctx context.Context, userID string) (*domain.DTOProfileResponse, error) {
	viewer, err := q.repository.GetByID(ctx, userID)
	if err != nil {
		return nil, base.ToDomainError(err)
	}

	return domain.NewDTOProfileResponse(viewer), nil
}
