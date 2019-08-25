package main

import (
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/pkg/log"
	"github.com/chromz/replicator/pkg/watcher"
	"os"
)

func main() {
	log.Info("Initializing replicator client...")
	err := config.LoadConfig()
	if err != nil {
		os.Exit(1)
	}
	directory := config.Directory()
	log.InitMessage(
		"rclient",
		"directory \""+directory+"\"",
	)
	watcher.Start(directory)
}
