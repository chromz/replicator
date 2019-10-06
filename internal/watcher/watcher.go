package watcher

import (
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/pkg/log"
	"github.com/fsnotify/fsnotify"
)

// EventFunction is a function to handle an inotify event
type EventFunction func(*fsnotify.Event)

var (
	onCreate EventFunction
	onWrite  EventFunction
	onRemove EventFunction
)

func doSync(event fsnotify.Event) {
	if event.Op&fsnotify.Create == fsnotify.Create {
		onCreate(&event)
	} else if event.Op&fsnotify.Write == fsnotify.Write {
		onWrite(&event)
	} else if event.Op&fsnotify.Remove == fsnotify.Remove {
		onRemove(&event)
	}
}

func watchFile(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			doSync(event)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Error("Inotify error", err)
		}
	}
}

// Start is a function to initialize watching for a file or directory for
// inotify events
func Start(createHandler, writeHandler, removeHandler EventFunction) {
	onCreate = createHandler
	onWrite = writeHandler
	onRemove = removeHandler
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("Unable to create watcher", err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go watchFile(watcher)
	directory := config.Directory()
	err = watcher.Add(directory)
	if err != nil {
		log.Error("Unable to add directory to watch list", err)
	}
	log.InitMessage(
		"watcher",
		"directory \""+directory+"\"",
	)
	<-done
}
