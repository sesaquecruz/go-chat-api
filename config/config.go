package config

import (
	"strings"

	"github.com/sesaquecruz/go-chat-api/pkg/log"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type BrokerConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type APIConfig struct {
	Port         string
	Path         string
	Mode         string
	AllowOrigins string
	JwtIssuer    string
	JwtAudience  string
}

type Config struct {
	Database DatabaseConfig
	Broker   BrokerConfig
	API      APIConfig
}

var env *viper.Viper
var file *viper.Viper
var cfg *Config

func init() {
	logger := log.NewLogger("config")

	env = viper.New()
	env.SetDefault("APP_DATABASE_HOST", "")
	env.SetDefault("APP_DATABASE_PORT", "")
	env.SetDefault("APP_DATABASE_USER", "")
	env.SetDefault("APP_DATABASE_PASSWORD", "")
	env.SetDefault("APP_DATABASE_NAME", "")
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

	err := file.ReadInConfig()
	if err != nil {
		logger.Info(err)
	}
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

	cfg.Broker = BrokerConfig{
		Host:     getString("APP_BROKER_HOST"),
		Port:     getString("APP_BROKER_PORT"),
		User:     getString("APP_BROKER_USER"),
		Password: getString("APP_BROKER_PASSWORD"),
	}

	cfg.API = APIConfig{
		Port:         getString("APP_API_PORT"),
		Path:         getString("APP_API_PATH"),
		Mode:         getString("APP_API_MODE"),
		AllowOrigins: getString("APP_API_CORS_ORIGINS"),
		JwtIssuer:    getString("APP_API_JWT_ISSUER"),
		JwtAudience:  getString("APP_API_JWT_AUDIENCE"),
	}
}

func GetConfig() Config {
	return *cfg
}
