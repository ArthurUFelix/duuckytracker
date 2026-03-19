package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIBaseURL  string
	APIUsername string
	APIPassword string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

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
