package config

import (
	"flag"
)

var (
	accessTokenSecretVal string
	refreshTokenSecretVal string
)

var (
	AccessTokenSecret = &accessTokenSecretVal
	RefreshTokenSecret = &refreshTokenSecretVal
)

func init() {
	if flag.Lookup("access-token-secret") == nil {
		flag.StringVar(&accessTokenSecretVal, "access-token-secret", "your-access-token-secret-change-in-production", "Access token secret")
	}
	if flag.Lookup("refresh-token-secret") == nil {
		flag.StringVar(&refreshTokenSecretVal, "refresh-token-secret", "your-refresh-token-secret-change-in-production", "Refresh token secret")
	}
}

func LoadConfig() *Config {
	return &Config{
		Auth: AuthConfig{
			AccessTokenSecret:  "your-access-token-secret-change-in-production",
			RefreshTokenSecret: "your-refresh-token-secret-change-in-production",
		},
	}
}