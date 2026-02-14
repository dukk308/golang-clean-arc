package config

import (
	"flag"
)

var (
	accessTokenSecretVal  string
	refreshTokenSecretVal string

	// for google oauth
	googleOauthClientID     string
	googleOauthClientSecret string
	googleOauthRedirectURL  string
	googleOauthAuthURL      string
	googleOauthTokenURL     string
	googleOauthUserInfoURL  string
	googleOauthScopes       string
)

var (
	AccessTokenSecret  = &accessTokenSecretVal
	RefreshTokenSecret = &refreshTokenSecretVal
)

func init() {
	if flag.Lookup("access-token-secret") == nil {
		flag.StringVar(&accessTokenSecretVal, "access-token-secret", "", "Access token secret")
	}
	if flag.Lookup("refresh-token-secret") == nil {
		flag.StringVar(&refreshTokenSecretVal, "refresh-token-secret", "", "Refresh token secret")
	}
	if flag.Lookup("google-client-id") == nil {
		flag.StringVar(&googleOauthClientID, "google-client-id", "", "Google client id")
	}
	if flag.Lookup("google-client-secret") == nil {
		flag.StringVar(&googleOauthClientSecret, "google-client-secret", "", "Google client secret")
	}
	if flag.Lookup("google-redirect-url") == nil {
		flag.StringVar(&googleOauthRedirectURL, "google-redirect-url", "", "Google redirect URL")
	}
	if flag.Lookup("google-auth-url") == nil {
		flag.StringVar(&googleOauthAuthURL, "google-auth-url", "", "Google auth URL")
	}
	if flag.Lookup("google-token-url") == nil {
		flag.StringVar(&googleOauthTokenURL, "google-token-url", "", "Google token URL")
	}
	if flag.Lookup("google-user-info-url") == nil {
		flag.StringVar(&googleOauthUserInfoURL, "google-user-info-url", "", "Google user info URL")
	}
	if flag.Lookup("google-scopes") == nil {
		flag.StringVar(&googleOauthScopes, "google-scopes", "", "Google scopes (comma-separated)")
	}
}

func LoadConfig() *Config {
	return &Config{
		Auth: AuthConfig{
			AccessTokenSecret:  accessTokenSecretVal,
			RefreshTokenSecret: refreshTokenSecretVal,
			GoogleClientID:     googleOauthClientID,
			GoogleClientSecret: googleOauthClientSecret,
			GoogleRedirectURL:  googleOauthRedirectURL,
			GoogleAuthURL:      googleOauthAuthURL,
			GoogleTokenURL:     googleOauthTokenURL,
			GoogleUserInfoURL:  googleOauthUserInfoURL,
			GoogleScopes:       parseScopes(googleOauthScopes),
		},
	}
}

func parseScopes(scopes string) []string {
	if scopes == "" {
		return nil
	}
	var result []string
	for _, s := range splitScopes(scopes) {
		if s := trimSpace(s); s != "" {
			result = append(result, s)
		}
	}
	return result
}

func splitScopes(scopes string) []string {
	var result []string
	start := 0
	for i := 0; i < len(scopes); i++ {
		if scopes[i] == ',' {
			result = append(result, scopes[start:i])
			start = i + 1
		}
	}
	result = append(result, scopes[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
