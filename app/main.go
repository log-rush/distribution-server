package main

import (
	"log"

	_lsHttpHandler "github.com/fabiankachlock/log-rush-simple-server/logstream/delivery/http"
	_lsUseCase "github.com/fabiankachlock/log-rush-simple-server/logstream/usecase"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	logStreamUseCase := _lsUseCase.NewLogStreamUSeCase()
	_lsHttpHandler.NewLogStreamHandler(app, logStreamUseCase)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Send([]byte("pong"))
	})

	log.Fatal(app.Listen("127.0.0.1:7000"))
}
