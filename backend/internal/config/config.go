package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port           string `envconfig:"PORT" default:"9000"`
	DatabaseURL    string `envconfig:"DATABASE_URL" required:"true"`
	JWTSecret      string `envconfig:"JWT_SECRET" required:"true" default:"your_very_strong_encypted_secret"`
	JWTExpiry      int    `envconfig:"JWT_EXPIRY" default:"24"` // in hours
	ApiVersion     string `envconfig:"API_VERSION" default:"v1.0.0"`
	PlunkBaseUrl   string `envconfig:"PLUNK_BASE_URL"`
	PlunkSecretKey string `envconfig:"PLUNK_SECRET_KEY"`
	RedisHost      string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPassword  string `envconfig:"REDIS_PASSWORD"`
	RedisPort      string `envconfig:"REDIS_PORT" default:"6379"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
