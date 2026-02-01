package domain

import "time"

type DTOCreateBlog struct {
	Title   string `json:"title" binding:"required"`
	Slug    string `json:"slug" binding:"required"`
	Content string `json:"content"`
}

type DTOBlogResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func NewDTOBlogResponse(blog *Blog) *DTOBlogResponse {
	return &DTOBlogResponse{
		ID:        blog.ID.String(),
		Title:     blog.Title,
		Slug:      blog.Slug,
		Content:   blog.Content,
		CreatedAt: blog.CreatedAt,
		UpdatedAt: blog.UpdatedAt,
	}
}

type DTOBlogListResponse struct {
	Items []*DTOBlogResponse `json:"items"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}
