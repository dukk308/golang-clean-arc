package repository

import (
	"context"

	"github.com/dukk308/golang-clean-arch-starter/internal/modules/auth/domain"
	user_domain "github.com/dukk308/golang-clean-arch-starter/internal/modules/user/domain"
)

type UserRepositoryAdapter struct {
	userRepository user_domain.IViewerRepository
}

func NewUserRepositoryAdapter(userRepository user_domain.IViewerRepository) domain.IUserRepository {
	return &UserRepositoryAdapter{
		userRepository: userRepository,
	}
}

func (a *UserRepositoryAdapter) GetByEmail(ctx context.Context, email string) (*domain.UserInfo, error) {
	viewer, err := a.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		ID:       viewer.ID.String(),
		Email:    viewer.Email.Value,
		Password: viewer.Password,
		Role:     viewer.Role.String(),
	}, nil
}

func (a *UserRepositoryAdapter) GetByID(ctx context.Context, id string) (*domain.UserInfo, error) {
	viewer, err := a.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		ID:       viewer.ID.String(),
		Email:    viewer.Email.Value,
		Password: viewer.Password,
		Role:     viewer.Role.String(),
	}, nil
}
