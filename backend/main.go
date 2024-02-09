package main

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db := connectDB()
	app := fiber.New()

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

	// users
	app.Get("/users", func(c *fiber.Ctx) error {
		// get parameter from query string, currPage
		currPageStr := c.Query("currPage")
		currPage, err := strconv.ParseUint(currPageStr, 10, 16)
		if err != nil {
			currPage = 1
		}
		return usersList(c, db, uint16(currPage))
	})

	app.Get("/get_msgs", isValidJwt, func(c *fiber.Ctx) error {
		return getMessages(c, db)
	})

	app.Post("/send_msg", isValidJwt, func(c *fiber.Ctx) error {
		return sendMessage(c, db)
	})

	app.Listen(":3001")
}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:930adf4c5b27ddde79e6@tcp(127.0.0.1:3307)/realtime_chat")

	if err != nil {
		panic(err.Error())
	}
	return db
}

var jwtSecret = []byte("secret")
