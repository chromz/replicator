package watcher

import (
	"bytes"
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/pkg/log"
	"github.com/fsnotify/fsnotify"
	"os/exec"
)

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
	host := config.RsyncServer().Address
	module := config.Module()
	url := host + "::" + module
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
		if dir[len(dir)-1:] != "/" {
			dir += "/"
		}
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

	done := make(chan bool)
	go watchFile(watcher)
	directory := config.Directory()
	err = watcher.Add(directory)
	if err != nil {
		log.Error("Unable to add directory to watch list", err)
	}
	log.InitMessage(
		"rclient",
		"directory \""+directory+"\"",
	)
	<-done
}
