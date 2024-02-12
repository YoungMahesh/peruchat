package main

import (
	"database/sql"
	"encoding/json"
)

type SendMessage struct {
	To      string `json:"to"`
	Message string `json:"message"`
}
type SendMessageEvent struct {
	Type    string       `json:"type"`
	Payload []GetMessage `json:"payload"`
}

type NewMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

type SendUsersEvent struct {
	Type    string `json:"type"`
	Payload []User `json:"payload"`
}

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func sendMessage0(username string, to string, message string, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO messages (from_user, to_user, message) VALUES (?, ?, ?)", username, to, message)
	if err != nil {
		return err
	}
	return nil
}

type GetMessagesRequest struct {
	ToUser string `json:"to_user"`
}

type GetMessage struct {
	IsSender bool   `json:"is_sender"`
	Message  string `json:"message"`
}
type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
}

func getMessages0(username string, to_user string, db *sql.DB) ([]GetMessage, error) {
	rows, err := db.Query("SELECT from_user, to_user, message FROM messages WHERE (from_user = ? AND to_user = ?) OR (to_user = ? AND from_user = ?)", username, to_user, username, to_user)
	if err != nil {
		return nil, err
	}

	var messages []GetMessage
	for rows.Next() {
		var m1 Message
		if err := rows.Scan(&m1.From, &m1.To, &m1.Message); err != nil {
			return nil, err
		}
		isSender := false
		if m1.From == username {
			isSender = true
		}
		gm1 := GetMessage{
			IsSender: isSender,
			Message:  m1.Message,
		}
		if m1.To == to_user {
			messages = append(messages, gm1)
		}
	}

	return messages, nil
}
