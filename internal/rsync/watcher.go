package rsync

import (
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/pkg/log"
	"github.com/fsnotify/fsnotify"
	"sync"
)

// EventQueue is a fifo struct of fsnotify events
type EventQueue struct {
	// Events is the slice of fsnotify events
	Events []fsnotify.Event
	// Mux that locks the struct
	Mux sync.Mutex
}

func watchFile(watcher *fsnotify.Watcher, eventQueue *EventQueue) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			eventQueue.Mux.Lock()
			eventQueue.Events = append(eventQueue.Events, event)
			eventQueue.Mux.Unlock()
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Error("Inotify error", err)
		}
	}
}

// Start is a function to start listening for inotify events
func Start() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("Unable to create watcher", err)
	}
	defer watcher.Close()
	host := config.RsyncServer().Address
	module := config.Module()
	url := host + "::" + module
	dir := config.Directory()

	eventQueue := &EventQueue{}
	done := make(chan bool)
	go watchFile(watcher, eventQueue)
	err = watcher.Add(dir)
	if err != nil {
		log.Error("Unable to add directory to watch list", err)
	}
	ticker := NewTicker(url, eventQueue)
	go ticker.Run()
	log.InitMessage(
		"rclient",
		"directory \""+dir+"\"",
	)
	<-done
}
