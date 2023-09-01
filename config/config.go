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

type CacheConfig struct {
	Host string
	Port string
}

type BrokerConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type ApiConfig struct {
	Port         string
	Path         string
	Mode         string
	AllowOrigins string
	JwtIssuer    string
	JwtAudience  string
}

type Config struct {
	Database DatabaseConfig
	Cache    CacheConfig
	Broker   BrokerConfig
	Api      ApiConfig
}

var (
	env  *viper.Viper
	file *viper.Viper
	cfg  *Config
)

func init() {
	env = viper.New()
	env.SetDefault("APP_DATABASE_HOST", "")
	env.SetDefault("APP_DATABASE_PORT", "")
	env.SetDefault("APP_DATABASE_USER", "")
	env.SetDefault("APP_DATABASE_PASSWORD", "")
	env.SetDefault("APP_DATABASE_NAME", "")
	env.SetDefault("APP_CACHE_HOST", "")
	env.SetDefault("APP_CACHE_PORT", "")
	env.SetDefault("APP_BROKER_HOST", "")
	env.SetDefault("APP_BROKER_PORT", "")
	env.SetDefault("APP_BROKER_USER", "")
	env.SetDefault("APP_BROKER_PASSWORD", "")
	env.SetDefault("APP_API_PORT", "")
	env.SetDefault("APP_API_PATH", "")
	env.SetDefault("APP_API_MODE", "")
	env.SetDefault("APP_API_CORS_ORIGINS", "")
	env.SetDefault("APP_API_JWT_ISSUER", "")
	env.SetDefault("APP_API_JWT_AUDIENCE", "")
	env.AutomaticEnv()

	file = viper.New()
	file.SetConfigName("config")
	file.SetConfigType("toml")
	file.AddConfigPath(".")
	file.ReadInConfig()
}

func getValue(key string) string {
	value := env.GetString(key)
	if value == "" {
		value = file.GetString(strings.Replace(key, "_", ".", -1))
	}

	return value
}

func Load() Config {
	cfg = new(Config)

	cfg.Database = DatabaseConfig{
		Host:     getValue("APP_DATABASE_HOST"),
		Port:     getValue("APP_DATABASE_PORT"),
		User:     getValue("APP_DATABASE_USER"),
		Password: getValue("APP_DATABASE_PASSWORD"),
		Name:     getValue("APP_DATABASE_NAME"),
	}

	cfg.Cache = CacheConfig{
		Host: getValue("APP_CACHE_HOST"),
		Port: getValue("APP_CACHE_PORT"),
	}

	cfg.Broker = BrokerConfig{
		Host:     getValue("APP_BROKER_HOST"),
		Port:     getValue("APP_BROKER_PORT"),
		User:     getValue("APP_BROKER_USER"),
		Password: getValue("APP_BROKER_PASSWORD"),
	}

	cfg.Api = ApiConfig{
		Port:         getValue("APP_API_PORT"),
		Path:         getValue("APP_API_PATH"),
		Mode:         getValue("APP_API_MODE"),
		AllowOrigins: getValue("APP_API_CORS_ORIGINS"),
		JwtIssuer:    getValue("APP_API_JWT_ISSUER"),
		JwtAudience:  getValue("APP_API_JWT_AUDIENCE"),
	}

	return *cfg
}
