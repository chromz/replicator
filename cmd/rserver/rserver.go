package main

import (
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/internal/server"
	"github.com/chromz/replicator/internal/watcher"
	"github.com/chromz/replicator/pkg/log"
	"os"
)

func main() {
	log.Info("Initializing replicator server...")
	configPath := ""
	if len(os.Args) == 2 {
		configPath = os.Args[1]
	}
	err := config.LoadConfig(configPath)
	if err != nil {
		os.Exit(1)
	}

	hub := server.NewHub()
	go hub.Run()
	go server.Hear(hub)
	watcher.Start(server.BroadcastHandler, server.UpdateHandler,
		server.RemoveHandler)
}
