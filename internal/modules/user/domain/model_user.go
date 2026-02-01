package domain

import (
	common "github.com/dukk308/beetool.dev-go-starter/pkgs/base"
)

type AuthProvider string

const (
	AuthProviderLocal  AuthProvider = "local"
	AuthProviderGoogle AuthProvider = "google"
)

type User struct {
	common.BaseModel
	Username       string       `json:"username"`
	Email          *EmailVO     `json:"email"`
	Password       string       `json:"password"`
	Role           Role         `json:"role"`
	AuthProvider   AuthProvider `json:"auth_provider"`
	AuthProviderID *string      `json:"auth_provider_id"`
}

func NewUser(username string, email string, password string) (*User, error) {
	emailVo := NewEmailVO(email)
	if err := emailVo.Validate(); err != nil {
		return nil, err
	}

	return &User{
		BaseModel: *common.GenerateBaseModel(),
		Username:  username,
		Email:     emailVo,
		Password:  password,
	}, nil
}
