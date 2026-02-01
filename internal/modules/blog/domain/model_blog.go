package domain

import (
	common "github.com/dukk308/beetool.dev-go-starter/pkgs/base"
)

type Blog struct {
	common.BaseModel
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Content string `json:"content"`
}

func NewBlog(title, slug, content string) *Blog {
	return &Blog{
		BaseModel: *common.GenerateBaseModel(),
		Title:     title,
		Slug:      slug,
		Content:   content,
	}
}
