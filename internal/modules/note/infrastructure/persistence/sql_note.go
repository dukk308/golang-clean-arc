package persistence

import (
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/note/domain"
	common "github.com/dukk308/beetool.dev-go-starter/pkgs/base"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/components/gorm_comp"
	"github.com/google/uuid"
)

type SQLNote struct {
	gorm_comp.SQLModel
	Title   string `gorm:"column:title;type:varchar(255);not null"`
	Slug    string `gorm:"column:slug;type:varchar(255);uniqueIndex:uni_notes_slug;not null"`
	Content string `gorm:"column:content;type:text"`
}

func (n *SQLNote) TableName() string {
	return "notes"
}

func (n *SQLNote) ToDomain() *domain.Note {
	return &domain.Note{
		BaseModel: common.BaseModel{
			ID:        uuid.MustParse(n.ID),
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
			DeletedAt: n.DeletedAt,
		},
		Title:   n.Title,
		Slug:    n.Slug,
		Content: n.Content,
	}
}

func (n *SQLNote) FromDomain(note *domain.Note) {
	n.ID = note.ID.String()
	n.CreatedAt = note.CreatedAt
	n.UpdatedAt = note.UpdatedAt
	n.DeletedAt = note.DeletedAt
	n.Title = note.Title
	n.Slug = note.Slug
	n.Content = note.Content
}
