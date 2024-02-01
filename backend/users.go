package main

import (
	"database/sql"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	Username string `json:"username"`
}

func usersList(c *fiber.Ctx, db *sql.DB, currPage uint16) error {

	usersPerPage := uint16(10)
	rows, err := db.Query("SELECT username FROM users LIMIT ? OFFSET ?", usersPerPage, (currPage-1)*usersPerPage)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not connect database to get users",
		})
	}

	var users []User
	for rows.Next() {
		var u1 User
		if err := rows.Scan(&u1.Username); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not extract users from data",
			})
		}
		users = append(users, u1)
	}

	return c.Status(http.StatusOK).JSON(users)
}
