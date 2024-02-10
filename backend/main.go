package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gofiber/contrib/websocket"
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

	app.Use("/ws", func(c *fiber.Ctx) error {
		println("someone sent a request to /ws")
		if websocket.IsWebSocketUpgrade(c) {
			println("upgrading to websocket")
			username, err := getUsernameFromToken(c.Query("token"))
			println("websocket-username", username)
			if err != nil || username == "" {
				return fiber.ErrUnauthorized
			}
			c.Locals("username", username)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		username := c.Locals("username").(string)
		println("username inside ws", username)
		// c.Locals is added to the *websocket.Conn
		// log.Println(c.Locals("allowed"))  // true
		// log.Println(c.Params("id"))       // 123
		// log.Println(c.Query("v"))         // 1.0
		// log.Println(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)

		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			println("got message", string(msg))
			var message1 Event
			err := json.Unmarshal(msg, &message1)
			if err != nil {
				log.Println("failed json.Unmarshal:", err)
				break
			}
			log.Printf("recv: %s", msg)
			if message1.Type == "get_msgs" {
				var getMsgReq GetMessagesRequest
				err := json.Unmarshal(message1.Payload, &getMsgReq)
				if err != nil {
					log.Println("failed json.Unmarshal:", err)
					break
				}

				messages, err := getMessages0(username, getMsgReq.ToUser, db)
				if err != nil {
					log.Println("failed getMessages0:", err)
					break
				}

				var sendMessages Event
				sendMessages.Type = "get_msgs_resp"
				sendMessages.Payload, err = json.Marshal(messages)
				if err != nil {
					log.Println("failed json.Marshal sendMessages:", err)
					break
				}

				c.WriteMessage(mt, sendMessages.Payload)

				println("get_msgs", messages)

			} else if message1.Type == "send_msg" {
				var sendMsg SendMessage
				err := json.Unmarshal(message1.Payload, &sendMsg)
				if err != nil {
					log.Println("failed json.Unmarshal:", err)
					break
				}

				err = sendMessage0(username, sendMsg.To, sendMsg.Message, db)
				if err != nil {
					log.Println("failed sendMessage0:", err)
					break
				}

				println("send_msg", message1.Payload)
			}

			// if err = c.WriteMessage(mt, msg); err != nil {
			// 	log.Println("write:", err)
			// 	break
			// }
		}

	}))

	app.Listen(":3001")
}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:930adf4c5b27ddde79e6@tcp(127.0.0.1:3307)/realtime_chat")

	if err != nil {
		panic(err.Error())
	}
	return db
}
