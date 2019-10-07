package client

import (
	"bytes"
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/pkg/log"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"net/url"
	"os/exec"
	"strings"
)

var rsyncUrl string

// Load loads all necessary information
func Load() {
	host := config.RsyncServer().Address
	module := config.Module()
	rsyncUrl = host + "::" + module
}

func ListenCreatedFiles(conn *websocket.Conn) {
	destDir := config.Directory()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("Unexpected: ", err)
			break
		}
		messageString := string(message)
		fileSlice := strings.Split(messageString, "$")
		switch fileSlice[0] {
		case "create":
			break
		default:
			break
		}
		fileName := "/" + string(fileSlice[1])
		rsyncHandler("Pulling new changes: ", "-avOzh",
			rsyncUrl+string(fileName), destDir)
	}
}

func ServerConn() *websocket.Conn {
	addr := url.URL{Scheme: "ws", Host: config.Host(), Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(addr.String(), nil)
	if err != nil {
		log.Info("Unable to connect to socket: ", err)
		return nil
	}
	return conn
}

func runRsync(params ...string) (string, string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("rsync", params...)
	log.Info(params)
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
	rsyncHandler("Created file: "+event.Name, "-avzh", event.Name, rsyncUrl)
}

// WriteHandler the handler when a file is updated
func WriteHandler(event *fsnotify.Event) {
	rsyncHandler("Updated file: "+event.Name, "-auvzh", event.Name, rsyncUrl)
}

// RemoveHandler the handler when a file is removed
func RemoveHandler(event *fsnotify.Event) {
	dir := config.Directory()
	rsyncHandler("Removed file: "+event.Name, "-avhO", "--delete", dir, rsyncUrl)
}
