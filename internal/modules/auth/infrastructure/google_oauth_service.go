package infrastructure

import (
	"context"
	"encoding/json"
	"io"

	"github.com/dukk308/beetool.dev-go-starter/internal/config"
	"github.com/dukk308/beetool.dev-go-starter/internal/modules/auth/domain"
	"golang.org/x/oauth2"
)

type GoogleOAuthService struct {
	config *oauth2.Config

	googleUserInfoURL string
}

func NewGoogleOAuthService(cfg *config.AuthConfig) domain.IGoogleOAuthService {

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes:       cfg.GoogleScopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.GoogleAuthURL,
			TokenURL: cfg.GoogleTokenURL,
		},
	}

	return &GoogleOAuthService{config: oauth2Config, googleUserInfoURL: cfg.GoogleUserInfoURL}
}

func (s *GoogleOAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *GoogleOAuthService) ExchangeCode(ctx context.Context, code string) (*domain.DTOGoogleUser, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := s.config.Client(ctx, token)
	resp, err := client.Get(s.googleUserInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var googleUser domain.DTOGoogleUser
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, err
	}

	return &googleUser, nil
}
