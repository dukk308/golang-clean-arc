package application

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
)

type ListBlogsQuery struct {
	repository domain.IBlogRepository
}

func NewListBlogsQuery(repository domain.IBlogRepository) *ListBlogsQuery {
	return &ListBlogsQuery{
		repository: repository,
	}
}

func (q *ListBlogsQuery) Execute(ctx context.Context, page, limit int) (*domain.DTOBlogListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	blogs, total, err := q.repository.GetPage(ctx, offset, limit)
	if err != nil {
		return nil, base.ToDomainError(err)
	}
	items := make([]*domain.DTOBlogResponse, len(blogs))
	for i, b := range blogs {
		items[i] = domain.NewDTOBlogResponse(b)
	}
	return &domain.DTOBlogListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}
