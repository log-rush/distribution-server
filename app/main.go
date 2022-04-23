package main

import (
	"log"

	_ "github.com/fabiankachlock/log-rush-simple-server/docs"
	_lsHttpHandler "github.com/fabiankachlock/log-rush-simple-server/logstream/delivery/http"
	_lsUseCase "github.com/fabiankachlock/log-rush-simple-server/logstream/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title log-rush simple server
// @version 1.0
// @description this is an example log-rush server
// @contact.name API Support
// @contact.email contact@fabiankachlock.dev
// @license.name n.a.
// @license.url n.a.
// @host localhost:7000
// @BasePath /
// @Tag.name logstream
// @Tag.description all endpoint for logstreams
func main() {
	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	logStreamUseCase := _lsUseCase.NewLogStreamUSeCase()
	_lsHttpHandler.NewLogStreamHandler(app, logStreamUseCase)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Send([]byte("pong"))
	})

	log.Fatal(app.Listen("127.0.0.1:7000"))
}
