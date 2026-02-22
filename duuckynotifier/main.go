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

	rabbitConsumer, err := consumer.NewConsumer(cfg.RabbitMQURL, "summoners", discord)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer rabbitConsumer.Close()

	err = rabbitConsumer.Start()
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	log.Println("Listening for messages... Press Ctrl+C to stop")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
}