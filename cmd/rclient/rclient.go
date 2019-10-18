package main

import (
	"flag"
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/internal/rsync"
	"github.com/chromz/replicator/pkg/log"
	"os"
)

func main() {
	log.Info("Initializing replicator client...")
	configPath := flag.String("c", "", "Config file path")
	flag.Parse()
	err := config.LoadConfig(configPath)
	if err != nil {
		os.Exit(1)
	}
	rsync.Start()
}
