package main

import (
	"context"
	"database/sql"
	"sync"
)

type Manager struct {
	clients  ClientsMap
	handlers map[string]EventHandler
	sync.RWMutex
	db *sql.DB
}

func NewManager(ctx context.Context, db *sql.DB) *Manager {
	m := &Manager{
		clients:  make(ClientsMap),
		handlers: make(map[string]EventHandler),
		db:       db,
	}

	m.Lock() // can use because of sync.RWMutex
	m.handlers["get_users"] = getUsersHandler
	m.handlers["get_msgs"] = getChatMessagesHandler
	m.handlers["send_msg"] = sendChatMessagesHandler
	m.Unlock()

	return m
}

func (manager *Manager) removeClient(client *Client) {
	manager.Lock()
	defer manager.Unlock()

	if _, ok := manager.clients[client]; ok {
		client.connection.Close()
		delete(manager.clients, client)
	}
}
