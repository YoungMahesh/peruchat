package main

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
)

type User struct {
	Username string `json:"username"`
}

func usersList0(db *sql.DB, currPage uint16) ([]User, error) {

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

func handleGetUsers(c *websocket.Conn, mt int, db *sql.DB) {
	users, err := usersList0(db, 1)
	if err != nil {
		log.Println("failed usersList0:", err)
		return
	}
	var sendUsersEvent SendUsersEvent
	sendUsersEvent.Type = "get_users_resp"
	sendUsersEvent.Payload = users
	sendUsersJson, err := json.Marshal(sendUsersEvent)
	if err != nil {
		log.Println("failed json.Marshal sendUsersEvent:", err)
		return
	}
	c.WriteMessage(mt, sendUsersJson)
	println("sent_users count", len(users))
}
