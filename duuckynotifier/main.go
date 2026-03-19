package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arthurufelix/duuckynotifier/internal/config"
	"github.com/arthurufelix/duuckynotifier/internal/consumer"
	"github.com/arthurufelix/duuckynotifier/internal/notifier"
)

func main() {
	log.Println("Starting Duucky Notifier...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	discord := notifier.NewDiscordNotifier(cfg.DiscordWebhookURL)

	summonerConsumer, err := consumer.NewConsumer(cfg.RabbitMQURL, "summoners", discord)
	if err != nil {
		log.Fatalf("Failed to create summoner consumer: %v", err)
	}
	defer summonerConsumer.Close()

	err = summonerConsumer.Start()
	if err != nil {
		log.Fatalf("Failed to start summoner consumer: %v", err)
	}

	matchConsumer, err := consumer.NewConsumer(cfg.RabbitMQURL, "matches", discord)
	if err != nil {
		log.Fatalf("Failed to create match consumer: %v", err)
	}
	defer matchConsumer.Close()

	err = matchConsumer.Start()
	if err != nil {
		log.Fatalf("Failed to start match consumer: %v", err)
	}

	log.Println("Listening for messages... Press Ctrl+C to stop")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
}
