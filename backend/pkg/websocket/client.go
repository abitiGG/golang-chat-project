package websocket

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
	mu   sync.Mutex
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		_ = c.Conn.Close()
	}()

	for {

		messageType, payload, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break // Exit the loop on error
		}

		msg := Message{Type: messageType, Body: string(payload)}
		c.Pool.Broadcast <- msg
		fmt.Printf("Received Message: %+v\n", msg)
	}
}
