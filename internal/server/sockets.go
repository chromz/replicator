package server

import (
	"github.com/chromz/replicator/internal/config"
	"github.com/chromz/replicator/pkg/log"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{}

var hub *Hub

func readMessages() {
}

func writeMessages() {
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Unable to create upgrader: ", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte)}
	hub.register <- client
}

// Hear starts socket server
func Hear(initHub *Hub) {
	hub = initHub
	http.HandleFunc("/", socketHandler)
	err := http.ListenAndServe(":"+config.Port(), nil)
	if err != nil {
		log.Error("Unable to start server: ", err)
	}
}
