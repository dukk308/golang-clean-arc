package config

type AuthConfig struct {
	AccessTokenSecret  string `mapstructure:"access_token_secret"`
	RefreshTokenSecret string `mapstructure:"refresh_token_secret"`
}

type Config struct {
	Auth     AuthConfig
}
