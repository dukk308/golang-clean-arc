package application

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/note/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
)

type UpdateNoteCommand struct {
	repository domain.INoteRepository
}

func NewUpdateNoteCommand(repository domain.INoteRepository) *UpdateNoteCommand {
	return &UpdateNoteCommand{
		repository: repository,
	}
}

func (c *UpdateNoteCommand) Execute(ctx context.Context, id string, dto *domain.DTOCreateNote) (*domain.DTONoteResponse, error) {
	note, err := c.repository.GetByID(ctx, id)
	if err != nil {
		return nil, base.ToDomainError(err)
	}
	note.Title = dto.Title
	note.Slug = dto.Slug
	note.Content = dto.Content
	if err := c.repository.Update(ctx, note); err != nil {
		return nil, base.ToDomainError(err)
	}
	return domain.NewDTONoteResponse(note), nil
}
