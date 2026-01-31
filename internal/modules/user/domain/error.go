package domain

import "github.com/dukk308/golang-clean-arch-starter/pkgs/ddd"

var (
	ErrInvalidUsername = &ddd.DomainError{
		Message: "username cannot be empty",
		Code:    "INVALID_USERNAME",
	}
	ErrInvalidEmail = &ddd.DomainError{
		Message: "email cannot be empty",
		Code:    "INVALID_EMAIL",
	}
	ErrInvalidRole = &ddd.DomainError{
		Message: "invalid role",
		Code:    "INVALID_ROLE",
	}
	ErrUnauthorized = &ddd.DomainError{
		Message: "unauthorized action",
		Code:    "UNAUTHORIZED",
	}
)
