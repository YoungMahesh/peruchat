package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
)

type ClientsMap map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	egress     chan Event
	username   string
}

func (manager *Manager) setupClient(conn *websocket.Conn, username string) *Client {
	client := &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
		username:   username,
	}

	manager.Lock()
	manager.clients[client] = true
	manager.Unlock()

	go client.readMessages()
	// TODO: make client.writeMessages() a goroutine
	// currently, if we make it goroutine, the program is closing websocket connection
	//   as goroutine works in parallel
	client.writeMessaages()

	println("--------------client setup done")
	return client
}

func (client *Client) readMessages() {
	var (
		// mt  int
		msg []byte
		err error
	)
	for {
		_, msg, err = client.connection.ReadMessage()
		if err != nil {
			println("client.readMessage: failed", err)
			break
		}

		println("client.readMessage got message", string(msg))
		var message1 Event
		err = json.Unmarshal(msg, &message1)
		if err != nil {
			println("readMessage: failed json.Unmarshal:", err)
			break
		}

		handler, ok := client.manager.handlers[message1.Type]
		if !ok {
			println("readMessage: handler not found for", message1.Type)
			continue
		}

		handler(client, message1)
	}
}

func (client *Client) writeMessaages() {
	for {
		message, ok := <-client.egress
		log.Println("client.writeMessages: received", message, ok)
		if !ok {
			println("client.writeMessages: egress not ok")
			continue
		}

		messageText, err := json.Marshal(message)
		if err != nil {
			log.Println("client.writeMessage: failed json.Marshal:", err)
			continue
		}

		println("sent_users", messageText)
		err = client.connection.WriteMessage(websocket.TextMessage, messageText)
		if err != nil {
			log.Println("failed client.writeMessage:", err)
		}
	}
}
