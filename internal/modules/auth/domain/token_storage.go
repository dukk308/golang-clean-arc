package domain

import "context"

type ITokenStorage interface {
	StoreRefreshToken(ctx context.Context, userID, token string) error
	GetRefreshToken(ctx context.Context, userID string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID string) error
	IsRefreshTokenValid(ctx context.Context, userID, token string) (bool, error)
}

type InMemoryTokenStorage struct {
	tokens map[string]string
}

func NewInMemoryTokenStorage() ITokenStorage {
	return &InMemoryTokenStorage{
		tokens: make(map[string]string),
	}
}

func (s *InMemoryTokenStorage) StoreRefreshToken(ctx context.Context, userID, token string) error {
	s.tokens[userID] = token
	return nil
}

func (s *InMemoryTokenStorage) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	token, exists := s.tokens[userID]
	if !exists {
		return "", ErrInvalidToken
	}
	return token, nil
}

func (s *InMemoryTokenStorage) DeleteRefreshToken(ctx context.Context, userID string) error {
	delete(s.tokens, userID)
	return nil
}

func (s *InMemoryTokenStorage) IsRefreshTokenValid(ctx context.Context, userID, token string) (bool, error) {
	storedToken, err := s.GetRefreshToken(ctx, userID)
	if err != nil {
		return false, err
	}
	return storedToken == token, nil
}
