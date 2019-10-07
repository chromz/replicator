package main

import (
	"github.com/chromz/replicator/internal/client"
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/internal/watcher"
	"github.com/chromz/replicator/pkg/log"
	"os"
)

func main() {
	log.Info("Initializing replicator client...")
	configPath := ""
	if len(os.Args) == 2 {
		configPath = os.Args[1]
	}
	err := config.LoadConfig(configPath)
	if err != nil {
		os.Exit(1)
	}
	conn := client.ServerConn()
	if conn == nil {
		os.Exit(1)
	}
	client.Load()
	go client.ListenCreatedFiles(conn)
	watcher.Start(client.CreateHandler, client.WriteHandler,
		client.RemoveHandler)
}
