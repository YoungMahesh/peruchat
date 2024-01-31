package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
)

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var jwtSecret = []byte("secret")

func userLogin(c *fiber.Ctx, db *sql.DB) error {
	u1 := new(Login)
	if err := c.BodyParser(u1); err != nil {
		return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{
			"message": "Cannot parse JSON",
		})
	}

	var dbUsername, dbPassword string
	err := db.QueryRow("SELECT username, password FROM users WHERE username = ?", u1.Username).Scan(&dbUsername, &dbPassword)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Username not found",
		})
	}

	// compare hashed password from user with hashed password from database
	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(u1.Password))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": dbUsername,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // token expires in 24 hours
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate token",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func userRegister(c *fiber.Ctx, db *sql.DB) error {
	u1 := new(Register)
	if err := c.BodyParser(u1); err != nil {
		return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{
			"message": "Cannot parse JSON",
		})
	}

	if !isValidUsername(u1.Username) {
		return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{
			"message": "Username should be 3 to 20 characters long and only contain alphanumeric characters",
		})
	}

	// check if username is already taken
	var dbUsername string
	err := db.QueryRow("SELECT username FROM users WHERE username = ?", u1.Username).Scan(&dbUsername)
	if err != sql.ErrNoRows {
		return c.Status(http.StatusNotAcceptable).JSON(fiber.Map{
			"message": "Username is already taken",
		})
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(u1.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not hash password",
		})
	}
	u1.Password = string(hashedPass)

	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", u1.Username, u1.Email, u1.Password)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not insert user to database",
		})
	}

	return c.Status(http.StatusCreated).JSON(u1)
}
