package domain

import (
	common "github.com/dukk308/beetool.dev-go-starter/pkgs/base"
)

type Note struct {
	common.BaseModel
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Content string `json:"content"`
}

func NewNote(title, slug, content string) *Note {
	return &Note{
		BaseModel: *common.GenerateBaseModel(),
		Title:     title,
		Slug:      slug,
		Content:   content,
	}
}
