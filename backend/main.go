package main

import (
	"context"
	"database/sql"
	// "encoding/json"
	// "log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use("/ws", func(c *fiber.Ctx) error {
		// println("someone sent a request to /ws")
		if websocket.IsWebSocketUpgrade(c) {
			// println("upgrading to websocket")
			username, err := getUsernameFromToken(c.Query("token"))
			// println("websocket-username", username)
			if err != nil || username == "" {
				println("username ")
				return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
			}
			c.Locals("username", username)
			return c.Next()
		}
		return c.Status(fiber.StatusUpgradeRequired).SendString("Upgrade Required")
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		username := c.Locals("username").(string)
		println("username inside ws", username)
		// c.Locals is added to the *websocket.Conn
		// log.Println(c.Locals("allowed"))  // true
		// log.Println(c.Params("id"))       // 123
		// log.Println(c.Query("v")) // 1.0
		// log.Println(c.Cookies("session")) // ""

		manager.setupClient(c, username)

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		// var (
		// 	mt  int
		// 	msg []byte
		// 	err error
		// )

		// for {
		// 	if mt, msg, err = c.ReadMessage(); err != nil {
		// 		log.Println("read:", err)
		// 		break
		// 	}
		// 	// println("got message", string(msg))
		// 	var message1 Event
		// 	err := json.Unmarshal(msg, &message1)
		// 	if err != nil {
		// 		log.Println("failed json.Unmarshal:", err)
		// 		break
		// 	}
		// 	// log.Printf("recv: %s", msg)
		// 	// if message1.Type == "get_users" {
		// 	// 	handleGetUsers(c, mt, db)
		// 	// } else
		// 	if message1.Type == "get_msgs" {
		// 		handleGetMsgs(c, mt, message1, username, db)
		// 	} else if message1.Type == "send_msg" {
		// 		handleSendMsg(c, message1, username, db)
		// 	}

		// 	// if err = c.WriteMessage(mt, msg); err != nil {
		// 	// 	log.Println("write:", err)
		// 	// 	break
		// 	// }
		// }

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
