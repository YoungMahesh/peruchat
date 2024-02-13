package main

import (
	"context"
	"sync"
)

type Manager struct {
	clients  ClientsMap
	handlers map[string]EventHandler
	sync.RWMutex
}

func NewManager(ctx context.Context) *Manager {
	m := &Manager{
		clients:  make(ClientsMap),
		handlers: make(map[string]EventHandler),
	}

	m.Lock() // can use because of sync.RWMutex
	m.handlers["get_users"] = getUsersHandler
	m.Unlock()

	return m
}
