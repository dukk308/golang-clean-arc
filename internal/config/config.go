package config

type AuthConfig struct {
	AccessTokenSecret  string   `mapstructure:"access_token_secret"`
	RefreshTokenSecret string   `mapstructure:"refresh_token_secret"`
	GoogleClientID     string   `mapstructure:"google_client_id"`
	GoogleClientSecret string   `mapstructure:"google_client_secret"`
	GoogleRedirectURL  string   `mapstructure:"google_redirect_url"`
	GoogleAuthURL      string   `mapstructure:"google_auth_url"`
	GoogleTokenURL     string   `mapstructure:"google_token_url"`
	GoogleUserInfoURL  string   `mapstructure:"google_user_info_url"`
	GoogleScopes       []string `mapstructure:"google_scopes"`
}

type Config struct {
	Auth AuthConfig
}
