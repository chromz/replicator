package server

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func NewClient() *Client {
	return &Client{
		send: make(chan []byte),
	}
}
