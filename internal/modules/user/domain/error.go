package domain

import "errors"

var (
	ErrInvalidUsername = errors.New("username cannot be empty")
	ErrInvalidEmail    = errors.New("email cannot be empty")
	ErrInvalidRole     = errors.New("invalid role")
	ErrUnauthorized    = errors.New("unauthorized action")
)
