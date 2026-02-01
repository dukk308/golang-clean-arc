package application

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/note/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
)

type ListNotesQuery struct {
	repository domain.INoteRepository
}

func NewListNotesQuery(repository domain.INoteRepository) *ListNotesQuery {
	return &ListNotesQuery{
		repository: repository,
	}
}

func (q *ListNotesQuery) Execute(ctx context.Context) ([]*domain.DTONoteResponse, error) {
	notes, err := q.repository.GetAll(ctx)
	if err != nil {
		return nil, base.ToDomainError(err)
	}
	result := make([]*domain.DTONoteResponse, len(notes))
	for i, n := range notes {
		result[i] = domain.NewDTONoteResponse(n)
	}
	return result, nil
}
