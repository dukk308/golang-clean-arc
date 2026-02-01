package application

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
)

type DeleteBlogCommand struct {
	repository domain.IBlogRepository
}

func NewDeleteBlogCommand(repository domain.IBlogRepository) *DeleteBlogCommand {
	return &DeleteBlogCommand{
		repository: repository,
	}
}

func (c *DeleteBlogCommand) Execute(ctx context.Context, id string) error {
	return base.ToDomainError(c.repository.Delete(ctx, id))
}
