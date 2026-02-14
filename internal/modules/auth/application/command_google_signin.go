package application

import (
	"context"

	"github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/domain"
	user_domain "github.com/dukk308/beetool.dev-go-starter/internal/modules/user/domain"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
	"gorm.io/gorm"
)

type GoogleSigninCommand struct {
	googleOAuthService domain.IGoogleOAuthService
	repository         domain.IUserRepository
	userRepository     user_domain.IViewerRepository
	tokenService       domain.ITokenService
	tokenStorage       domain.ITokenStorage
}

func NewGoogleSigninCommand(
	googleOAuthService domain.IGoogleOAuthService,
	repository domain.IUserRepository,
	userRepository user_domain.IViewerRepository,
	tokenService domain.ITokenService,
	tokenStorage domain.ITokenStorage,
) *GoogleSigninCommand {
	return &GoogleSigninCommand{
		googleOAuthService: googleOAuthService,
		repository:         repository,
		userRepository:     userRepository,
		tokenService:       tokenService,
		tokenStorage:       tokenStorage,
	}
}

func (c *GoogleSigninCommand) Execute(ctx context.Context, dto *domain.DTOGoogleSignin) (*domain.DTOTokenResponse, error) {
	googleUser, err := c.googleOAuthService.ExchangeCode(ctx, dto.Code)
	if err != nil {
		return nil, base.ToDomainError(err)
	}

	user, err := c.repository.GetByEmail(ctx, googleUser.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new user with Google auth
			user, err = c.createGoogleUser(ctx, googleUser)
			if err != nil {
				return nil, base.ToDomainError(err)
			}
		} else {
			return nil, base.ToDomainError(err)
		}
	}

	accessToken, err := c.tokenService.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, base.ToDomainError(err)
	}

	refreshToken, err := c.tokenService.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, base.ToDomainError(err)
	}

	if err := c.tokenStorage.StoreRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, base.ToDomainError(err)
	}

	return &domain.DTOTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (c *GoogleSigninCommand) createGoogleUser(ctx context.Context, googleUser *domain.DTOGoogleUser) (*domain.UserInfo, error) {
	userID := googleUser.ID
	userDTO := &user_domain.DTOCreateUser{
		Username: googleUser.Name,
		Email:    googleUser.Email,
		Password: "", // No password for Google users
		Provider: user_domain.AuthProviderGoogle,
	}

	viewer, err := user_domain.CreateViewer(userDTO)
	if err != nil {
		return nil, err
	}

	// Set the AuthProviderID to the Google user ID
	viewer.User.AuthProviderID = &userID

	if err := c.userRepository.Create(ctx, viewer); err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		ID:       viewer.ID.String(),
		Email:    viewer.Email.Value,
		Password: viewer.Password,
		Role:     viewer.Role.String(),
	}, nil
}
