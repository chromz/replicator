package rsync

import (
	"bytes"
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/pkg/log"
	"github.com/fsnotify/fsnotify"
	"os/exec"
)

var url string

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

func doSync(event fsnotify.Event) {
	if event.Op&fsnotify.Create == fsnotify.Create {
		stdout, stderr, err := runRsync("-avzh", event.Name, url)
		if err != nil {
			log.Error(stderr, err)
			return
		}
		log.Info("Created file: ", event.Name)
		log.Info(stdout)
	} else if event.Op&fsnotify.Write == fsnotify.Write {
		stdout, stderr, err := runRsync("-auvzh", event.Name, url)
		if err != nil {
			log.Error(stderr, err)
			return
		}
		log.Info("Updated file: ", event.Name)
		log.Info(stdout)
	} else if event.Op&fsnotify.Remove == fsnotify.Remove {
		dir := config.Directory()
		stdout, stderr, err := runRsync("-avhO", "--delete", dir, url)
		if err != nil {
			log.Error(stderr, err)
			return
		}
		log.Info("Deleted file: ", event.Name)
		log.Info(stdout)
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

func Start() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("Unable to create watcher", err)
	}
	defer watcher.Close()
	host := config.RsyncServer().Address
	module := config.Module()
	url = host + "::" + module
	done := make(chan bool)
	go watchFile(watcher)
	directory := config.Directory()
	err = watcher.Add(directory)
	if err != nil {
		log.Error("Unable to add directory to watch list", err)
	}
	ticker := NewTicker(directory, url, config.PollingRate())
	go ticker.Run()
	log.InitMessage(
		"rclient",
		"directory \""+directory+"\"",
	)
	<-done
}
