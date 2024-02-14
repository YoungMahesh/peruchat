package main

import (
	"encoding/json"
	"log"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(client *Client, event Event) error

func getUsersHandler(client *Client, event Event) error {
	// get data from database
	currPage := uint16(1)
	usersPerPage := uint16(10)
	rows, err := client.manager.db.Query("SELECT username FROM users LIMIT ? OFFSET ?", usersPerPage, (currPage-1)*usersPerPage)
	if err != nil {
		return err
	}
	type User struct {
		Username string `json:"username"`
	}
	var users []User
	for rows.Next() {
		var u1 User
		if err := rows.Scan(&u1.Username); err != nil {
			return err
		}
		users = append(users, u1)
	}

	// send data to client
	var sendUsersEvent Event
	sendUsersEvent.Type = "get_users_resp"
	sendUsersEvent.Payload, err = json.Marshal(users)
	if err != nil {
		log.Println("failed json.Marshal sendUsersEvent:", err)
		return err
	}
	client.egress <- sendUsersEvent
	return nil
}

func getChatMessagesHandler(client *Client, event Event) error {
	println("----------------------------- getChatMessagesHandler --------------------------")

	// extract payload from event
	type GetMessagesRequest struct {
		ToUser string `json:"to_user"`
	}
	var getMsgReq GetMessagesRequest
	err := json.Unmarshal(event.Payload, &getMsgReq)
	if err != nil {
		log.Println("failed json.Unmarshal:", err)
		return err
	}

	// get data from database
	rows, err := client.manager.db.Query("SELECT from_user, to_user, message FROM messages WHERE (from_user = ? AND to_user = ?) OR (to_user = ? AND from_user = ?)", client.username, getMsgReq.ToUser, client.username, getMsgReq.ToUser)
	if err != nil {
		println("getChatMessagesHandler: failed to fetch from database", err)
		return err
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
	var messages []GetMessage
	for rows.Next() {
		var m1 Message
		if err := rows.Scan(&m1.From, &m1.To, &m1.Message); err != nil {
			return err
		}
		isSender := false
		if m1.From == client.username {
			isSender = true
		}
		gm1 := GetMessage{
			IsSender: isSender,
			Message:  m1.Message,
		}
		messages = append(messages, gm1)
	}

	// send data to client
	var sendMessages Event
	sendMessages.Type = "get_msgs_resp"
	sendMessages.Payload, err = json.Marshal(messages)
	if err != nil {
		log.Println("getChatMessageHandler: failed json.Marshal:", err)
		return err
	}
	client.egress <- sendMessages
	return nil
}

func sendChatMessagesHandler(client *Client, event Event) error {
	println("----------------------------- sendChat --------------------------")

	// extract payload from event
	type SendMessage struct {
		To      string `json:"to"`
		Message string `json:"message"`
	}
	var sendMsg SendMessage
	err := json.Unmarshal(event.Payload, &sendMsg)
	if err != nil {
		log.Println("failed json.Unmarshal:", err)
		return err
	}
	// send data to database
	_, err = client.manager.db.Exec("INSERT INTO messages (from_user, to_user, message) VALUES (?, ?, ?)", client.username, sendMsg.To, sendMsg.Message)
	if err != nil {
		log.Println("sendChatMessagesHandler: failed to insert into database:", err)
		return err
	}

	// send data to client
	var storedMessage Event
	storedMessage.Type = "send_msg_resp"
	type StoredMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Message string `json:"message"`
	}
	storedMessage.Payload, err = json.Marshal(&StoredMessage{
		From:    client.username,
		To:      sendMsg.To,
		Message: sendMsg.Message,
	})
	if err != nil {
		log.Println("storedMessageHandler: failed json.Marshal:", err)
		return err
	}
	client.egress <- storedMessage
	return nil
}
