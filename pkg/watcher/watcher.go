package watcher

import (
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/pkg/log"
	"github.com/fsnotify/fsnotify"
)

func watchFile(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Info(event)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Error("Inotify error", err)
		}
	}
}

func Start(directory string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("Unable to create watcher", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go watchFile(watcher)
	err = watcher.Add(config.GetDirectory())
	if err != nil {
		log.Error("Unable to add directory to watch list", err)
	}
	<-done
}
