package application

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/note/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
)

type GetNoteQuery struct {
	repository domain.INoteRepository
}

func NewGetNoteQuery(repository domain.INoteRepository) *GetNoteQuery {
	return &GetNoteQuery{
		repository: repository,
	}
}

func (q *GetNoteQuery) ExecuteByID(ctx context.Context, id string) (*domain.DTONoteResponse, error) {
	note, err := q.repository.GetByID(ctx, id)
	if err != nil {
		return nil, base.ToDomainError(err)
	}
	return domain.NewDTONoteResponse(note), nil
}

func (q *GetNoteQuery) ExecuteBySlug(ctx context.Context, slug string) (*domain.DTONoteResponse, error) {
	note, err := q.repository.GetBySlug(ctx, slug)
	if err != nil {
		return nil, base.ToDomainError(err)
	}
	return domain.NewDTONoteResponse(note), nil
}
