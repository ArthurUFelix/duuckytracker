package config

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

//go:embed .env
var embeddedEnv string

type Config struct {
	APIBaseURL  string
	APIUsername string
	APIPassword string
}

func loadEmbeddedConfig() error {
	lines := strings.Split(strings.TrimSpace(embeddedEnv), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			err := os.Setenv(key, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Load() (*Config, error) {
	err := loadEmbeddedConfig()
	if err != nil {
		_ = godotenv.Load()
	}

	apiBaseUrl := os.Getenv("API_BASE_URL")
	apiUsername := os.Getenv("API_USERNAME")
	apiPassword := os.Getenv("API_PASSWORD")

	if apiBaseUrl == "" {
		return nil, fmt.Errorf("API_BASE_URL is required")
	}
	if apiUsername == "" {
		return nil, fmt.Errorf("API_USERNAME is required")
	}
	if apiPassword == "" {
		return nil, fmt.Errorf("API_PASSWORD is required")
	}

	return &Config{
		APIBaseURL:  apiBaseUrl,
		APIUsername: apiUsername,
		APIPassword: apiPassword,
	}, nil
}
