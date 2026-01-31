package application

import (
	"context"

	"github.com/dukk308/golang-clean-arch-starter/internal/modules/auth/domain"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/ddd"
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
		return ddd.ToDomainError(err)
	}

	if err := c.tokenStorage.DeleteRefreshToken(ctx, claims.UserID); err != nil {
		return ddd.ToDomainError(err)
	}

	return nil
}
