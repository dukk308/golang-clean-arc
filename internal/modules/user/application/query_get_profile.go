package application

import (
	"context"

	"github.com/dukk308/golang-clean-arch-starter/internal/modules/user/domain"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/ddd"
)

type GetProfileQuery struct {
	repository domain.IViewerRepository
}

func NewGetProfileQuery(repository domain.IViewerRepository) *GetProfileQuery {
	return &GetProfileQuery{
		repository: repository,
	}
}

func (q *GetProfileQuery) Execute(ctx context.Context, userID string) (*domain.DTOProfileResponse, error) {
	viewer, err := q.repository.GetByID(ctx, userID)
	if err != nil {
		return nil, ddd.ToDomainError(err)
	}

	return domain.NewDTOProfileResponse(viewer), nil
}
