package rsync

import (
	"bytes"
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/pkg/log"
	"github.com/fsnotify/fsnotify"
	"os/exec"
	"time"
)

// Synchronizer is a struct that ticks every polling rate
type Synchronizer struct {
	ticker *time.Ticker
}

var rsyncURL string
var destDir string
var tempDir string
var eventQueue *EventQueue

// NewTicker constructor of the synchronizer ticker
func NewTicker(url string, queue *EventQueue) *Synchronizer {
	rsyncURL = url
	destDir = config.Directory()
	tempDir = config.TempDir()
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
			_, stderr, err := runRsync("-avzhP", event.Name, rsyncURL)
			if err != nil {
				log.Error(stderr, err)
				newQueue = append(newQueue, event)
				continue
			}
			log.Info("Created file: ", event.Name)
		} else if event.Op&fsnotify.Write == fsnotify.Write {
			_, stderr, err := runRsync("-auvzhP", event.Name, rsyncURL)
			if err != nil {
				log.Error(stderr, err)
				newQueue = append(newQueue, event)
				continue
			}
			log.Info("Updated file: ", event.Name)
		} else if event.Op&fsnotify.Remove == fsnotify.Remove {
			dir := config.Directory()
			_, stderr, err := runRsync("-avhOP", dir, rsyncURL, "--delete")
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
	stdout, stderr, err := runRsync("-avuOzhP", "-T", tempDir,
		rsyncURL, destDir, "--delete")
	if err != nil {
		log.Error("Unable to run rsync pull: "+stderr, err)
		return
	}
	log.Info(stdout)
}

// Run starts the ticker
func (synchronizer *Synchronizer) Run() {
	log.Info("Trying to do initial synchronization")
	doneInitialSync := false
	for !doneInitialSync {
		doneInitialSync = true
		_, stderr, err := runRsync("-avuOzhP", destDir, rsyncURL)
		if err != nil {
			log.Error(stderr, err)
			doneInitialSync = false
		}
		time.Sleep(time.Millisecond * time.Duration(config.PollingRate()))
	}
	log.Info("Initial synchronization done")
	for {
		select {
		case <-synchronizer.ticker.C:
			sweepQueue()
			pullChanges()
		}
	}
}
