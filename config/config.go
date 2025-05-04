package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func Load[T any]() (T, error) {
	var cfg T

	if envfile := os.Getenv("ENV_FILE"); envfile != "" {
		if err := godotenv.Load(envfile); err != nil {
			return cfg, fmt.Errorf("failed to load env file: %w", err)
		}
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, fmt.Errorf("failed to process env variables: %w", err)
	}

	return cfg, nil
}
