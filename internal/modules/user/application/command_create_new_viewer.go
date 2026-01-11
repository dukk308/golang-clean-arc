package application

import (
	"context"

	"github.com/dukk308/golang-clean-arc/internal/modules/user/domain"
	"github.com/dukk308/golang-clean-arc/pkgs/ddd"
)

type CreateNewViewerCommand struct {
	repository domain.IViewerRepository
}

func NewCreateNewViewerCommand(repository domain.IViewerRepository) *CreateNewViewerCommand {
	return &CreateNewViewerCommand{
		repository: repository,
	}
}

func (c *CreateNewViewerCommand) Execute(ctx context.Context, dto *domain.DTOCreateUser) (*ddd.BaseModel, error) {
	viewer, err := domain.CreateViewer(dto)
	if err != nil {
		return nil, ddd.ToDomainError(err)
	}

	if err := c.repository.Create(ctx, viewer); err != nil {
		return nil, ddd.ToDomainError(err)
	}

	return &viewer.BaseModel, nil
}
