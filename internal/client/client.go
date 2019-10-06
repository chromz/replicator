package client

import (
	"bytes"
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/pkg/log"
	"github.com/fsnotify/fsnotify"
	"os/exec"
)

var url string

// Load loads all necessary information
func Load() {
	host := config.RsyncServer().Address
	module := config.Module()
	url = host + "::" + module
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

func rsyncHandler(message string, command ...string) {
	stdout, stderr, err := runRsync(command...)
	if err != nil {
		log.Error(stderr, err)
		return
	}
	log.Info(message)
	log.Info(stdout)
}

// CreateHandler the handler for new files
func CreateHandler(event *fsnotify.Event) {
	rsyncHandler("Created file: "+event.Name, "-avzh", event.Name, url)
}

// WriteHandler the handler when a file is updated
func WriteHandler(event *fsnotify.Event) {
	rsyncHandler("Updated file: "+event.Name, "-auvzh", event.Name, url)
}

// RemoveHandler the handler when a file is removed
func RemoveHandler(event *fsnotify.Event) {
	dir := config.Directory()
	rsyncHandler("Removed file: "+event.Name, "-avhO", "--delete", dir, url)
}
