package main

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type SendMessage struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func sendMessage(c *fiber.Ctx, db *sql.DB) error {
	sm1 := new(SendMessage)
	if err := c.BodyParser(sm1); err != nil {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"message": "Cannot parse JSON",
		})
	}

	username := c.Locals("username").(string)
	_, err := db.Exec("INSERT INTO messages (from_user, to_user, message) VALUES (?, ?, ?)", username, sm1.To, sm1.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not send message",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Message sent",
	})
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

func getMessages(c *fiber.Ctx, db *sql.DB) error {
	username := c.Locals("username").(string)
	to_user := c.Query("to")
	if to_user == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "message_recipient username is not provided",
		})
	}

	rows, err := db.Query("SELECT from_user, to_user, message FROM messages WHERE (from_user = ? AND to_user = ?) OR (to_user = ? AND from_user = ?)", username, to_user, username, to_user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get messages",
		})
	}

	var messages []GetMessage
	for rows.Next() {
		var m1 Message
		if err := rows.Scan(&m1.From, &m1.To, &m1.Message); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not extract messages from data",
			})
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

	return c.Status(fiber.StatusOK).JSON(messages)
}
