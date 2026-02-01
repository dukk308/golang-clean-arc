package application

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
)

type UpdateBlogCommand struct {
	repository domain.IBlogRepository
}

func NewUpdateBlogCommand(repository domain.IBlogRepository) *UpdateBlogCommand {
	return &UpdateBlogCommand{
		repository: repository,
	}
}

func (c *UpdateBlogCommand) Execute(ctx context.Context, id string, dto *domain.DTOCreateBlog) (*domain.DTOBlogResponse, error) {
	blog, err := c.repository.GetByID(ctx, id)
	if err != nil {
		return nil, base.ToDomainError(err)
	}
	blog.Title = dto.Title
	blog.Slug = dto.Slug
	blog.Content = dto.Content
	if err := c.repository.Update(ctx, blog); err != nil {
		return nil, base.ToDomainError(err)
	}
	return domain.NewDTOBlogResponse(blog), nil
}
