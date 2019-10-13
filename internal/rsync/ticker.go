package rsync

import (
	"bytes"
	"github.com/chromz/replicator/pkg/log"
	"os/exec"
	"time"
)

type Synchronizer struct {
	ticker *time.Ticker
}

var rsyncUrl string
var destDir string

func NewTicker(directory, url string, pollingRate int) *Synchronizer {
	rsyncUrl = url
	destDir = directory
	return &Synchronizer{
		ticker: time.NewTicker(time.Millisecond * time.Duration(pollingRate)),
	}
}

func pullChanges() {
	var stdout, stderr bytes.Buffer
	log.Info("Pulling changes")
	cmd := exec.Command("rsync", "-avOzh", rsyncUrl, destDir)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Error("Unable to run rsync pull: "+stderr.String(), err)
		return
	}
	log.Info(stdout.String())
}

func (synchronizer *Synchronizer) Run() {
	for {
		select {
		case <-synchronizer.ticker.C:
			pullChanges()
		}
	}
}
