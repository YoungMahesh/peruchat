package main

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
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

func handleSendMsg(c *websocket.Conn, message1 Event, username string, db *sql.DB) {
	var sendMsg SendMessage
	err := json.Unmarshal(message1.Payload, &sendMsg)
	if err != nil {
		log.Println("failed json.Unmarshal:", err)
		return
	}

	err = sendMessage0(username, sendMsg.To, sendMsg.Message, db)
	if err != nil {
		log.Println("failed sendMessage0:", err)
		return
	}

	println("send_msg", message1.Payload)
}
