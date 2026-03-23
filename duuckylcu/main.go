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

	log.Println("authenticating...")
	err = apiClient.Login(cfg.APIUsername, cfg.APIPassword)
	if err != nil {
		log.Fatalf("failed to authenticate with API: %v", err)
	}

	var lockfile *lcu.LockfileData
	for {
		lockfile, err = lcu.ReadLockfile()
		if err != nil {
			log.Println("waiting for League Client to start...")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	log.Println("connected to League Client")

	wsClient := lcu.NewWebsocketClient(lockfile.Port, lockfile.Password)

	if err := wsClient.Connect(); err != nil {
		log.Fatalf("failed to connect to websocket: %v", err)
	}
	defer wsClient.Close()

	if err := wsClient.Subscribe("OnJsonApiEvent_lol-chat_v1_friends"); err != nil {
		log.Fatalf("failed to subscribe to events: %v", err)
	}

	friendTracker := tracker.NewFriendTracker(wsClient, apiClient)
	go friendTracker.Start()

	log.Println("tracker is running! Press Ctrl + C to stop")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down gracefully...")
}
