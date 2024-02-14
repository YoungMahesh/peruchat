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
	m.Unlock()

	return m
}
