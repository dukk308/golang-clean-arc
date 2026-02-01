package application

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
)

type SignoutCommand struct {
	tokenStorage domain.ITokenStorage
	tokenService domain.ITokenService
}

func NewSignoutCommand(
	tokenStorage domain.ITokenStorage,
	tokenService domain.ITokenService,
) *SignoutCommand {
	return &SignoutCommand{
		tokenStorage: tokenStorage,
		tokenService: tokenService,
	}
}

func (c *SignoutCommand) Execute(ctx context.Context, refreshToken string) error {
	claims, err := c.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return base.ToDomainError(err)
	}

	if err := c.tokenStorage.DeleteRefreshToken(ctx, claims.UserID); err != nil {
		return base.ToDomainError(err)
	}

	return nil
}
