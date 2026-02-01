package domain

import "github.com/dukk308/beetool.dev-go-starter/pkgs/base"

var (
	ErrInvalidUsername = &base.DomainError{
		Message: "username cannot be empty",
		Code:    "INVALID_USERNAME",
	}
	ErrInvalidEmail = &base.DomainError{
		Message: "email cannot be empty",
		Code:    "INVALID_EMAIL",
	}
	ErrInvalidRole = &base.DomainError{
		Message: "invalid role",
		Code:    "INVALID_ROLE",
	}
	ErrUnauthorized = &base.DomainError{
		Message: "unauthorized action",
		Code:    "UNAUTHORIZED",
	}
)
