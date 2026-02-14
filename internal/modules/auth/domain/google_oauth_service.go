package domain

import "context"

// IGoogleOAuthService interface for Google OAuth operations
type IGoogleOAuthService interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*DTOGoogleUser, error)
}
