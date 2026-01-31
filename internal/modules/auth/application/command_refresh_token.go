package application

import (
	"context"

	"github.com/dukk308/golang-clean-arch-starter/internal/modules/auth/domain"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/ddd"
)

type RefreshTokenCommand struct {
	repository   domain.IUserRepository
	tokenService domain.ITokenService
	tokenStorage domain.ITokenStorage
}

func NewRefreshTokenCommand(
	repository domain.IUserRepository,
	tokenService domain.ITokenService,
	tokenStorage domain.ITokenStorage,
) *RefreshTokenCommand {
	return &RefreshTokenCommand{
		repository:   repository,
		tokenService: tokenService,
		tokenStorage: tokenStorage,
	}
}

func (c *RefreshTokenCommand) Execute(ctx context.Context, dto *domain.DTORefreshToken) (*domain.DTOTokenResponse, error) {
	claims, err := c.tokenService.ValidateRefreshToken(dto.RefreshToken)
	if err != nil {
		return nil, ddd.ToDomainError(err)
	}

	isValid, err := c.tokenStorage.IsRefreshTokenValid(ctx, claims.UserID, dto.RefreshToken)
	if err != nil || !isValid {
		return nil, ddd.ToDomainError(domain.ErrInvalidToken)
	}

	user, err := c.repository.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ddd.ToDomainError(domain.ErrInvalidToken)
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
