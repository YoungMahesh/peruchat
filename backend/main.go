package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	db := connectDB()
	app := fiber.New()

	rootContext := context.Background()
	ctx, cancel := context.WithCancel(rootContext)
	defer cancel()
	manager := NewManager(ctx, db)

	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	// authentication
	app.Post("/login", func(c *fiber.Ctx) error {
		return userLogin(c, db)
	})
	app.Post("/register", func(c *fiber.Ctx) error {
		return userRegister(c, db)
	})

	app.Get("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			username, err := getUsernameFromToken(c.Query("token"))
			println("token", c.Query("token"))
			if err != nil || username == "" {
				println("username ")
				return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
			}

			// verify if username exists in the database
			var dbUsername string
			err = db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&dbUsername)
			if err != nil {
				return c.Status(fiber.StatusNotFound).SendString("Username not found")
			}

			println("username inside ws", username)

			return websocket.New(func(wconn *websocket.Conn) {
				manager.setupClient(wconn, username)
			})(c)
		}
		return c.Status(fiber.StatusUpgradeRequired).SendString("Upgrade Required")
	})

	app.Listen(":3001")
}

func connectDB() *sql.DB {
	err := godotenv.Load("../database/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	password := os.Getenv("MYSQL_ROOT_PASSWORD")
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(127.0.0.1:3308)/peru_db", password))

	if err != nil {
		panic(err.Error())
	}
	return db
}
