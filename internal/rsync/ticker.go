package rsync

import (
	"bytes"
	"github.com/chromz/replicator/pkg/log"
	"os/exec"
	"time"
)

// Synchronizer is a struct that ticks every polling rate
type Synchronizer struct {
	ticker *time.Ticker
}

var rsyncURL string
var destDir string

// NewTicker constructor of the synchronizer ticker
func NewTicker(directory, url string, pollingRate int) *Synchronizer {
	rsyncURL = url
	destDir = directory
	return &Synchronizer{
		ticker: time.NewTicker(time.Millisecond * time.Duration(pollingRate)),
	}
}

func pullChanges() {
	var stdout, stderr bytes.Buffer
	log.Info("Pulling changes")
	cmd := exec.Command("rsync", "-avOzh", rsyncURL, destDir, "--delete")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Error("Unable to run rsync pull: "+stderr.String(), err)
		return
	}
	log.Info(stdout.String())
}

// Run starts the ticker
func (synchronizer *Synchronizer) Run() {
	for {
		select {
		case <-synchronizer.ticker.C:
			pullChanges()
		}
	}
}
