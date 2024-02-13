package main

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
)

type ClientsMap map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	egress     chan Event
	chatroom   string
}

func (manager *Manager) CreateClient(c *websocket.Conn) *Client {
	client := &Client{
		connection: c,
		manager:    manager,
		egress:     make(chan Event),
	}

	manager.clients[client] = true
	return client
}

func (client *Client) readMessages() {

	for {
		mt, msg, err := client.connection.ReadMessage()
		if err != nil {
			println("readMessage: failed", err)
			break
		}

		var message1 Event
		err = json.Unmarshal(msg, &message1)
		if err != nil {
			println("readMessage: failed json.Unmarshal:", err)
			break
		}

		handler, ok := client.manager.handlers[message1.Type]
		if !ok {
			println("readMessage: handler not found")
			break
		}

		users, err := handler(client, message1)

	}
}
