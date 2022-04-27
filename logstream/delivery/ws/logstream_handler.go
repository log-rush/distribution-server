package ws

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/log-rush/simple-server/domain"
)

type logStreamWsHandler struct {
	conns map[string]*websocket.Conn
	cu    domain.ClientsUseCase
	l     *domain.Logger
}

func NewLogStreamWsHandler(app *fiber.App, clientManager domain.ClientsUseCase, logger domain.Logger) {
	handler := &logStreamWsHandler{
		conns: map[string]*websocket.Conn{},
		cu:    clientManager,
		l:     &logger,
	}

	app.Use("/subscribe", handler.AllowWebsocketUpgrades)
	app.Get("/subscribe", websocket.New(handler.Connect))
}

func (h *logStreamWsHandler) AllowWebsocketUpgrades(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (h *logStreamWsHandler) Connect(conn *websocket.Conn) {
	var (
		mt  int = 1
		msg []byte
		err error
	)

	client, err := h.cu.NewClient(context.Background())
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
		h.cu.DestroyClient(context.Background(), client.ID)
		(*h.l).Debugf("closed connection %s", client.ID)
	}()

	go func() {
		for {
			select {
			case <-closed:
				return
			case message := <-client.Send:
				if message == nil {
					return
				}
				err = conn.WriteMessage(mt, message)
				if err != nil {
					(*h.l).Warnf("[%s] error while sending message: (%d) %s", err, mt, message)
				}
			}
		}
	}()

	(*h.l).Debugf("connected %s", client.ID)
	for {
		if mt, msg, err = conn.ReadMessage(); err != nil {
			if !websocket.IsCloseError(err, 1000) {
				(*h.l).Errorf("[%s] error while receiving message %s", client.ID, err)
			}
			break
		}
		client.Receive <- msg
	}

}
