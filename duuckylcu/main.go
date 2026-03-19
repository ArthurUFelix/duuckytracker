package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arthurufelix/duuckylcu/config"
	"github.com/arthurufelix/duuckylcu/internal/api"
	"github.com/arthurufelix/duuckylcu/internal/lcu"
	"github.com/arthurufelix/duuckylcu/internal/tracker"
)

func main() {
	log.Println("starting DuuckyLCU ARAM Tracker...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	log.Println("configuration loaded")

	apiClient := api.NewClient(cfg.APIBaseURL)

	log.Printf("authenticating as %s...", cfg.APIUsername)
	err = apiClient.Login(cfg.APIUsername, cfg.APIPassword)
	if err != nil {
		log.Fatalf("failed to authenticate with API: %v", err)
	}

	var lcuClient *lcu.Client
	for {
		lcuClient, err = lcu.NewClient()
		if err != nil {
			log.Println("waiting for League Client to start...")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	log.Println("connected to League Client")

	friendTracker := tracker.NewFriendTracker(lcuClient, apiClient)

	go friendTracker.Start(10 * time.Second)

	log.Println("tracker is running! Press Ctrl + C to stop")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down gracefully...")
}
