package rsync

import (
	"bytes"
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/pkg/log"
	"github.com/fsnotify/fsnotify"
	"os"
	"os/exec"
	"time"
)

// Synchronizer is a struct that ticks every polling rate
type Synchronizer struct {
	ticker *time.Ticker
}

var rsyncURL string
var tempDir string
var eventQueue *EventQueue
var watcher *fsnotify.Watcher

// NewTicker constructor of the synchronizer ticker
func NewTicker(url string, queue *EventQueue,
	fileWatcher *fsnotify.Watcher) *Synchronizer {
	rsyncURL = url
	tempDir = config.TempDir()
	watcher = fileWatcher
	eventQueue = queue
	return &Synchronizer{
		ticker: time.NewTicker(time.Millisecond *
			time.Duration(config.PollingRate())),
	}
}

func runRsync(params ...string) (string, string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("rsync", params...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", stderr.String(), err
	}
	return stdout.String(), "", nil
}

func sweepQueue() {
	eventQueue.Mux.Lock()
	log.Info("Sweeping event queue")
	newQueue := eventQueue.Events[:0]
	for _, event := range eventQueue.Events {
		if event.Op&fsnotify.Create == fsnotify.Create {
			if _, err := os.Stat(event.Name); os.IsNotExist(err) {
				log.Warn("Unable to create file " + event.Name +
					": File does not exist")
				continue
			}
			if fileInfo, err := os.Stat(event.Name); err == nil {
				mode := fileInfo.Mode()
				// Add a watcher if create file is directory
				if mode.IsDir() {
					log.Info("Watching directory: " +
						event.Name)
					watcher.Add(event.Name)
				}
			} else {
				log.Error("Could not create file", err)
				continue
			}
			_, stderr, err := runRsync("-avzhROP", event.Name, rsyncURL)
			if err != nil {
				log.Error(stderr, err)
				newQueue = append(newQueue, event)
				continue
			}
			log.Info("Created file: ", event.Name)
		} else if event.Op&fsnotify.Write == fsnotify.Write {
			if _, err := os.Stat(event.Name); os.IsNotExist(err) {
				log.Warn("Unable to create file " + event.Name +
					": File does not exist")
				continue
			}
			_, stderr, err := runRsync("-auvzhROP", event.Name, rsyncURL)
			if err != nil {
				log.Error(stderr, err)
				newQueue = append(newQueue, event)
				continue
			}
			log.Info("Updated file: ", event.Name)
		} else if event.Op&fsnotify.Remove == fsnotify.Remove {
			if _, err := os.Stat(event.Name); err == nil {
				err = os.Remove(event.Name)
				if err != nil {
					log.Error("Unable to remove file "+
						event.Name, err)
				}
				continue
			}
			if fileInfo, err := os.Stat(event.Name); err == nil {
				mode := fileInfo.Mode()
				// Add a watcher if create file is directory
				if mode.IsDir() {
					watcher.Remove(event.Name)
				}
			}
			_, stderr, err := runRsync("-avhORP", ".", rsyncURL, "--delete")
			if err != nil {
				log.Error(stderr, err)
				newQueue = append(newQueue, event)
				continue
			}
			log.Info("Deleted file: ", event.Name)

		}

	}
	eventQueue.Events = newQueue
	eventQueue.Mux.Unlock()
}

func pullChanges() {
	log.Info("Pulling changes")
	eventQueue.Mux.Lock()
	for _, event := range eventQueue.Events {
		if (event.Op&fsnotify.Create == fsnotify.Create) ||
			(event.Op&fsnotify.Remove == fsnotify.Remove) {
			log.Info("Waiting queue to clean deletes and creates...")
			eventQueue.Mux.Unlock()
			return
		}
	}
	eventQueue.Mux.Unlock()
	stdout, stderr, err := runRsync("-avOzhRP", "-T", tempDir,
		rsyncURL, ".", "--delete")
	if err != nil {
		log.Error("Unable to run rsync pull: "+stderr, err)
		return
	}
	log.Info(stdout)
}

// Run starts the ticker
func (synchronizer *Synchronizer) Run() {
	doneInitialSync := true
	if config.SyncOnStart() {
		doneInitialSync = false
	}
	for !doneInitialSync {
		log.Info("Trying to do initial synchronization")
		_, stderr, err := runRsync("-avuOzhRP", ".", rsyncURL)
		if err == nil {
			doneInitialSync = true
			log.Info("Initial synchronization done")
		} else {
			log.Error(stderr, err)
		}
		time.Sleep(time.Millisecond * time.Duration(config.PollingRate()))
	}
	for {
		select {
		case <-synchronizer.ticker.C:
			sweepQueue()
			pullChanges()
		}
	}
}
