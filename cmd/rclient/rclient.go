package main

import (
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/internal/watcher"
	"github.com/chromz/replicator/pkg/log"
	"os"
)

func main() {
	log.Info("Initializing replicator client...")
	err := config.LoadConfig()
	if err != nil {
		os.Exit(1)
	}
	watcher.Start()
}
