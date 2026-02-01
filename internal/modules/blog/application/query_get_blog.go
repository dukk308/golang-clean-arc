package application

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
)

type GetBlogQuery struct {
	repository domain.IBlogRepository
}

func NewGetBlogQuery(repository domain.IBlogRepository) *GetBlogQuery {
	return &GetBlogQuery{
		repository: repository,
	}
}

func (q *GetBlogQuery) ExecuteByID(ctx context.Context, id string) (*domain.DTOBlogResponse, error) {
	blog, err := q.repository.GetByID(ctx, id)
	if err != nil {
		return nil, base.ToDomainError(err)
	}
	return domain.NewDTOBlogResponse(blog), nil
}

func (q *GetBlogQuery) ExecuteBySlug(ctx context.Context, slug string) (*domain.DTOBlogResponse, error) {
	blog, err := q.repository.GetBySlug(ctx, slug)
	if err != nil {
		return nil, base.ToDomainError(err)
	}
	return domain.NewDTOBlogResponse(blog), nil
}
