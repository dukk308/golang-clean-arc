package persistence

import (
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/blog/domain"
	common "github.com/dukk308/beetool.dev-go-starter/pkgs/base"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/components/gorm_comp"
	"github.com/google/uuid"
)

type SQLBlog struct {
	gorm_comp.SQLModel
	Title   string `gorm:"column:title;type:varchar(255);not null"`
	Slug    string `gorm:"column:slug;type:varchar(255);uniqueIndex:uni_blogs_slug;not null"`
	Content string `gorm:"column:content;type:text"`
}

func (b *SQLBlog) TableName() string {
	return "blogs"
}

func (b *SQLBlog) ToDomain() *domain.Blog {
	return &domain.Blog{
		BaseModel: common.BaseModel{
			ID:        uuid.MustParse(b.ID),
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
			DeletedAt: b.DeletedAt,
		},
		Title:   b.Title,
		Slug:    b.Slug,
		Content: b.Content,
	}
}

func (b *SQLBlog) FromDomain(blog *domain.Blog) {
	b.ID = blog.ID.String()
	b.CreatedAt = blog.CreatedAt
	b.UpdatedAt = blog.UpdatedAt
	b.DeletedAt = blog.DeletedAt
	b.Title = blog.Title
	b.Slug = blog.Slug
	b.Content = blog.Content
}
