package domain

import (
	common "github.com/dukk308/golang-clean-arc/pkgs/ddd"
)

type User struct {
	common.BaseModel
	Username string   `json:"username"`
	Email    *EmailVO `json:"email"`
	Password string   `json:"password"`
	Role     Role     `json:"role"`
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
