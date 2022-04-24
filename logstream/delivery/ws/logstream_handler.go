package ws

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func NewLogStreamWsHandler(app *fiber.App) {
	app.Use("/subscribe", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/subscribe/:id", websocket.New(func(c *websocket.Conn) {
		var (
			mt  int
			msg []byte
			err error
		)

		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				fmt.Println("read:", err)
				break
			}
			fmt.Printf("recv: %s\n", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				fmt.Println("write:", err)
				break
			}
		}

	}))
}
