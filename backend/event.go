package main

import (
	"database/sql"
	"encoding/json"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(db *sql.DB, client *Client, event Event) ([]User, error)

func getUsersHandler(db *sql.DB, client *Client, event Event) ([]User, error) {
	currPage := uint16(1)
	usersPerPage := uint16(10)
	rows, err := db.Query("SELECT username FROM users LIMIT ? OFFSET ?", usersPerPage, (currPage-1)*usersPerPage)
	if err != nil {
		return nil, err
	}

	var users []User
	for rows.Next() {
		var u1 User
		if err := rows.Scan(&u1.Username); err != nil {
			return nil, err
		}
		users = append(users, u1)
	}

	return users, nil
}
