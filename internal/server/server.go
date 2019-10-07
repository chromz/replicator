package server

import (
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/pkg/log"
	"github.com/fsnotify/fsnotify"
	"strings"
)

var socketHub *Hub

func Load(hub *Hub) {
	socketHub = hub
}

func BroadcastHandler(event *fsnotify.Event) {
	dir := config.Directory()
	fileName := "create$" + strings.Replace(event.Name, dir, "", 1)
	hub.broadcast <- []byte(fileName)
	log.Info("Broadcasting file: " + fileName)
}

func UpdateHandler(event *fsnotify.Event) {
	dir := config.Directory()
	fileName := "update$" + strings.Replace(event.Name, dir, "", 1)
	hub.broadcast <- []byte(fileName)
	log.Info("Broadcasting file: " + fileName)
}

func RemoveHandler(event *fsnotify.Event) {
	dir := config.Directory()
	fileName := "remove$" + strings.Replace(event.Name, dir, "", 1)
	hub.broadcast <- []byte(fileName)
	log.Info("Broadcasting file: " + fileName)
}
