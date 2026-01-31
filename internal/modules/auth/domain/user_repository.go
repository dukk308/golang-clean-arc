package domain

import "context"

type IUserRepository interface {
	GetByEmail(ctx context.Context, email string) (*UserInfo, error)
	GetByID(ctx context.Context, id string) (*UserInfo, error)
}

type UserInfo struct {
	ID       string
	Email    string
	Password string
	Role     string
}
