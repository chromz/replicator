package server

import (
	"github.com/chromz/replicator/pkg/log"
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			log.Info("Incoming client")
			hub.clients[client] = true
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				close(client.send)
				delete(hub.clients, client)
			}
		case message := <-hub.broadcast:
			for client := range hub.clients {
				// Broadcast the message to all
				client.conn.WriteMessage(websocket.TextMessage,
					message)
			}
		}
	}
}
