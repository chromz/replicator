package main

import (
	"github.com/chromz/replicator/internal/config"
	// "github.com/chromz/replicator/internal/server"
	// "github.com/chromz/replicator/internal/watcher"
	"github.com/chromz/replicator/pkg/log"
	"os"
)

func main() {
	log.Info("Initializing replicator server...")
	err := config.LoadConfig()
	if err != nil {
		os.Exit(1)
	}
	log.Info("Sending created file")
	// server.Load()
	//watcher.Start(server.CreateHandler, server.WriteHandler,
	//server.RemoveHandler)
}
