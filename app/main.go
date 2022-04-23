package main

import (
	"log"
	"time"

	_ "github.com/fabiankachlock/log-rush-simple-server/docs"
	_lsHttpHandler "github.com/fabiankachlock/log-rush-simple-server/logstream/delivery/http"
	_lsRepo "github.com/fabiankachlock/log-rush-simple-server/logstream/repository/memory"
	_lsUseCase "github.com/fabiankachlock/log-rush-simple-server/logstream/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${method} - ${path} - ${status} (${latency}) \n",
	}))
	app.Use(cors.New())
	app.Use(recover.New())

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	logStreamRepo := _lsRepo.NewLogStreamRepository()
	logStreamUseCase := _lsUseCase.NewLogStreamUseCase(logStreamRepo, time.Second*3)
	_lsHttpHandler.NewLogStreamHandler(app, logStreamUseCase)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Send([]byte("pong"))
	})

	log.Fatal(app.Listen("127.0.0.1:7000"))
}
