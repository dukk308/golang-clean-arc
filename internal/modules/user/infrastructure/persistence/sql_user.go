package persistence

import (
	"github.com/dukk308/golang-clean-arch-starter/internal/modules/user/domain"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/components/gorm_comp"
	common "github.com/dukk308/golang-clean-arch-starter/pkgs/ddd"
	"github.com/google/uuid"
)

type SQLUser struct {
	gorm_comp.SQLModel
	AuthProvider   domain.AuthProvider `gorm:"column:auth_provider;type:varchar(255);default:local;index:idx_provider_id;index:idx_provider_username"`
	AuthProviderID *string             `gorm:"column:auth_provider_id;type:varchar(255);index:idx_provider_id"`
	Username       *string             `gorm:"column:username;type:varchar(255);index:idx_provider_username"`
	Email          *string             `gorm:"column:email;type:varchar(255)"`
	Password       *string             `gorm:"column:password;type:varchar(255)"`
	Role           domain.Role         `gorm:"column:role;type:varchar(255);default:viewer"`
}

func (u *SQLUser) TableName() string {
	return "users"
}

func (u *SQLUser) ToDomainViewer() *domain.Viewer {
	return &domain.Viewer{
		User: domain.User{
			BaseModel: common.BaseModel{
				ID:        uuid.MustParse(u.ID),
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
				DeletedAt: u.DeletedAt,
			},
			Username:       *u.Username,
			Email:          domain.NewEmailVO(*u.Email),
			Password:       *u.Password,
			Role:           domain.Role(u.Role),
			AuthProvider:   u.AuthProvider,
			AuthProviderID: u.AuthProviderID,
		},
	}
}

func (u *SQLUser) FromDomainViewer(viewer *domain.Viewer) {
	u.ID = viewer.ID.String()
	u.CreatedAt = viewer.CreatedAt
	u.UpdatedAt = viewer.UpdatedAt
	u.DeletedAt = viewer.DeletedAt
	u.AuthProvider = viewer.AuthProvider
	u.AuthProviderID = viewer.AuthProviderID
	u.Username = &viewer.Username
	u.Email = &viewer.Email.Value
	u.Password = &viewer.Password
	u.Role = viewer.Role
}
