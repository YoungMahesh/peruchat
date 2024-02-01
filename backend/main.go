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

	app.Listen(":3001")
}

func isValidUsername(str string) bool {
	// should be between 3 and 20 characters long
	if len(str) < 3 || len(str) > 20 {
		return false
	}
	// is alphanumeric
	for _, char := range str {
		if !(char >= 'a' && char <= 'z') && !(char >= '0' && char <= '9') {
			return false
		}
	}
	return true
}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:930adf4c5b27ddde79e6@tcp(127.0.0.1:3307)/realtime_chat")

	if err != nil {
		panic(err.Error())
	}
	return db
}

var jwtSecret = []byte("secret")
