package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"golang.org/x/crypto/bcrypt"
)

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	db := connectDB()
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

		if !isValidUsername(p1.Username) {
			return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{
				"message": "Username should be 3 to 20 characters long and only contain alphanumeric characters",
			})
		}

		// check if username is already taken
		var dbUsername string
		err := db.QueryRow("SELECT username FROM users WHERE username = ?", p1.Username).Scan(&dbUsername)
		if err != sql.ErrNoRows {
			return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{
				"message": "Username is already taken",
			})
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(p1.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not hash password",
			})
		}
		p1.Password = string(hashedPass)

		_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", p1.Username, p1.Email, p1.Password)
		if err != nil {
			log.Println(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not insert user to database",
			})
		}

		return c.Status(http.StatusCreated).JSON(p1)
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
