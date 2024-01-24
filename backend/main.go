package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		p1 := new(Register)
		if err := c.BodyParser(p1); err != nil {
			return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{
				"message": "Cannot parse JSON",
			})
		}
		return c.JSON(p1)
	})

	app.Listen(":3001")
}
