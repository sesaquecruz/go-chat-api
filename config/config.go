package config

import (
	"strings"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type APIConfig struct {
	Port         string
	Path         string
	GinMode      string
	AllowOrigins string
	JwtIssuer    string
	JwtAudience  string
}

type Config struct {
	API      APIConfig
	Database DatabaseConfig
}

var env *viper.Viper
var file *viper.Viper
var cfg *Config

func init() {
	env = viper.New()
	env.SetDefault("APP_DATABASE_HOST", "")
	env.SetDefault("APP_DATABASE_PORT", "")
	env.SetDefault("APP_DATABASE_USER", "")
	env.SetDefault("APP_DATABASE_PASSWORD", "")
	env.SetDefault("APP_DATABASE_NAME", "")
	env.SetDefault("APP_API_PORT", "")
	env.SetDefault("APP_API_PATH", "")
	env.SetDefault("APP_API_GIN_MODE", "")
	env.SetDefault("APP_API_ALLOW_ORIGINS", "")
	env.SetDefault("APP_API_JWT_ISSUER", "")
	env.SetDefault("APP_API_JWT_AUDIENCE", "")
	env.AutomaticEnv()

	file = viper.New()
	file.SetConfigName("config")
	file.SetConfigType("toml")
	file.AddConfigPath(".")
	file.ReadInConfig()
}

func getString(key string) string {
	value := env.GetString(key)
	if value == "" {
		value = file.GetString(strings.Replace(key, "_", ".", -1))
	}

	return value
}

func Load() {
	cfg = new(Config)

	cfg.Database = DatabaseConfig{
		Host:     getString("APP_DATABASE_HOST"),
		Port:     getString("APP_DATABASE_PORT"),
		User:     getString("APP_DATABASE_USER"),
		Password: getString("APP_DATABASE_PASSWORD"),
		Name:     getString("APP_DATABASE_NAME"),
	}

	cfg.API = APIConfig{
		Port:         getString("APP_API_PORT"),
		Path:         getString("APP_API_PATH"),
		GinMode:      getString("APP_API_GIN_MODE"),
		AllowOrigins: getString("APP_API_CORS_ORIGINS"),
		JwtIssuer:    getString("APP_API_JWT_ISSUER"),
		JwtAudience:  getString("APP_API_JWT_AUDIENCE"),
	}
}

func GetConfig() Config {
	return *cfg
}
