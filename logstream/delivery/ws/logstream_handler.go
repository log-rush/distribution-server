package ws

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/log-rush/simple-server/domain"
)

type logStreamWsHandler struct {
	conns         map[string]*websocket.Conn
	clientManager domain.ClientsUseCase
}

func NewLogStreamWsHandler(app *fiber.App, clientManager domain.ClientsUseCase) {
	handler := &logStreamWsHandler{
		conns:         map[string]*websocket.Conn{},
		clientManager: clientManager,
	}

	app.Use("/subscribe", handler.AllowWebsocketUpgrades)
	app.Get("/subscribe", func(ctx *fiber.Ctx) error {
		wsHandler := websocket.New(func(conn *websocket.Conn) {
			handler.Connect(conn, ctx)
		})
		return wsHandler(ctx)
	})
}

func (h *logStreamWsHandler) AllowWebsocketUpgrades(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (h *logStreamWsHandler) Connect(conn *websocket.Conn, ctx *fiber.Ctx) {
	var (
		mt  int
		msg []byte
		err error
	)

	client, err := h.clientManager.NewClient(ctx.Context())
	if err != nil {
		conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
		conn.Close()
		return
	}
	closed := make(chan bool)

	defer func() {
		closed <- true
		close(closed)
		conn.WriteMessage(websocket.CloseMessage, []byte{})
		conn.Close()
		h.clientManager.DestroyClient(ctx.Context(), client.ID)
	}()

	go func() {
		for {
			select {
			case <-closed:
				break
			case <-client.Send:
				if err = conn.WriteMessage(mt, msg); err != nil {
					break
				}
			}
		}
	}()

	for {
		if mt, msg, err = conn.ReadMessage(); err != nil {
			break
		}
		client.Receive <- msg
	}

}
