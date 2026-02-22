package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQURL string
	DiscordWebhookURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	if rabbitmqURL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL environment variable is required")
	}

	if webhookURL == "" {
		return nil, fmt.Errorf("DISCORD_WEBHOOK_URL environment variable is required")
	}

	return &Config{
		RabbitMQURL: rabbitmqURL,
		DiscordWebhookURL: webhookURL,
	}, nil
}