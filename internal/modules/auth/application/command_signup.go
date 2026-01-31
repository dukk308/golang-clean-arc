package application

import (
	"context"

	"github.com/dukk308/golang-clean-arch-starter/internal/modules/auth/domain"
	user_domain "github.com/dukk308/golang-clean-arch-starter/internal/modules/user/domain"
	"github.com/dukk308/golang-clean-arch-starter/pkgs/ddd"
)

type SignupCommand struct {
	repository   user_domain.IViewerRepository
	tokenService domain.ITokenService
}

func NewSignupCommand(
	repository user_domain.IViewerRepository,
	tokenService domain.ITokenService,
) *SignupCommand {
	return &SignupCommand{
		repository:   repository,
		tokenService: tokenService,
	}
}

func (c *SignupCommand) Execute(ctx context.Context, dto *domain.DTOSignup) (*ddd.BaseModel, error) {
	hashedPassword, err := c.tokenService.HashPassword(dto.Password)
	if err != nil {
		return nil, ddd.ToDomainError(err)
	}

	userDTO := &user_domain.DTOCreateUser{
		Username: dto.Username,
		Email:    dto.Email,
		Password: hashedPassword,
		Provider: user_domain.AuthProviderLocal,
	}

	viewer, err := user_domain.CreateViewer(userDTO)
	if err != nil {
		return nil, ddd.ToDomainError(err)
	}

	if err := c.repository.Create(ctx, viewer); err != nil {
		return nil, ddd.ToDomainError(err)
	}

	return &viewer.BaseModel, nil
}
