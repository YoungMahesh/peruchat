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

func handleGetMsgs(c *websocket.Conn, mt int, message1 Event, username string, db *sql.DB) {
	var getMsgReq GetMessagesRequest
	err := json.Unmarshal(message1.Payload, &getMsgReq)
	if err != nil {
		log.Println("failed json.Unmarshal:", err)
		return
	}

	messages, err := getMessages0(username, getMsgReq.ToUser, db)
	if err != nil {
		log.Println("failed getMessages0:", err)
		return
	}

	var sendMessages SendMessageEvent
	sendMessages.Type = "get_msgs_resp"
	sendMessages.Payload = messages

	sendMessagesJson, err := json.Marshal(sendMessages)
	if err != nil {
		log.Println("failed json.Marshal sendMessages:", err)
		return
	}
	c.WriteMessage(mt, sendMessagesJson)

	println("get_msgs", messages)
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
