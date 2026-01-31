package application

import (
	"context"

	"github.com/dukk308/golang-clean-arch-starter/internal/modules/auth/domain"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/ddd"
)

type SigninCommand struct {
	repository   domain.IUserRepository
	tokenService domain.ITokenService
	tokenStorage domain.ITokenStorage
}

func NewSigninCommand(
	repository domain.IUserRepository,
	tokenService domain.ITokenService,
	tokenStorage domain.ITokenStorage,
) *SigninCommand {
	return &SigninCommand{
		repository:   repository,
		tokenService: tokenService,
		tokenStorage: tokenStorage,
	}
}

func (c *SigninCommand) Execute(ctx context.Context, dto *domain.DTOSignin) (*domain.DTOTokenResponse, error) {
	user, err := c.repository.GetByEmail(ctx, dto.Email)
	if err != nil {
		return nil, ddd.ToDomainError(domain.ErrInvalidCredentials)
	}

	if err := c.tokenService.ComparePassword(user.Password, dto.Password); err != nil {
		return nil, ddd.ToDomainError(domain.ErrInvalidCredentials)
	}

	accessToken, err := c.tokenService.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, ddd.ToDomainError(err)
	}

	refreshToken, err := c.tokenService.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, ddd.ToDomainError(err)
	}

	if err := c.tokenStorage.StoreRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, ddd.ToDomainError(err)
	}

	return &domain.DTOTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
