package server

import (
	"github.com/chromz/replicator/pkg/log"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Unable to create upgrader: ", err)
		return
	}

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("Read error: ", err)
			continue
		}
	}

	conn.Close()

}

// Hear starts socket server
func Hear() {
	http.HandleFunc("/", socketHandler)
}
