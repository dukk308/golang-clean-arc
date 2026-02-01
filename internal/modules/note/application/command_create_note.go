package application

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/note/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
)

type CreateNoteCommand struct {
	repository domain.INoteRepository
}

func NewCreateNoteCommand(repository domain.INoteRepository) *CreateNoteCommand {
	return &CreateNoteCommand{
		repository: repository,
	}
}

func (c *CreateNoteCommand) Execute(ctx context.Context, dto *domain.DTOCreateNote) (*domain.DTONoteResponse, error) {
	note := domain.NewNote(dto.Title, dto.Slug, dto.Content)
	if err := c.repository.Create(ctx, note); err != nil {
		return nil, base.ToDomainError(err)
	}
	return domain.NewDTONoteResponse(note), nil
}
