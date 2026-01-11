package config

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type GlobalConfig struct {
	ServiceName string `mapstructure:"service_name"`
	Environment string `mapstructure:"environment"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type Config struct {
	Global   GlobalConfig
	Database DatabaseConfig
	Server   ServerConfig
}

func NewConfig() *Config {
	return &Config{
		Global: GlobalConfig{
			ServiceName: "golang-clean-arc",
			Environment: "development",
		},
		Server: ServerConfig{
			Port: "8080",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "3306",
			Username: "root",
			Password: "root",
			Database: "test",
		},
	}
}
