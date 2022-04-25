package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_ "github.com/log-rush/simple-server/docs"
	_lHttpHandler "github.com/log-rush/simple-server/log/delivery/http"
	_lRepo "github.com/log-rush/simple-server/log/repository/memory"
	_lUseCase "github.com/log-rush/simple-server/log/usecase"
	_lsHttpHandler "github.com/log-rush/simple-server/logstream/delivery/http"
	_lsWsHandler "github.com/log-rush/simple-server/logstream/delivery/ws"
	_lsRepo "github.com/log-rush/simple-server/logstream/repository/memory"
	_lsUseCase "github.com/log-rush/simple-server/logstream/usecase"
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
// @Tag.description all endpoints for logstreams
// @Tag.name log
// @Tag.description all endpoints for logs
func main() {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${method} - ${path} - ${status} (${latency}) \n",
	}))
	app.Use(cors.New())
	app.Use(recover.New())

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	logRepo := _lRepo.NewLogRepository(100)
	logStreamRepo := _lsRepo.NewLogStreamRepository()

	logStreamUseCase := _lsUseCase.NewLogStreamUseCase(logStreamRepo, time.Second*3)
	logUseCase := _lUseCase.NewLogUseCase(logRepo, logStreamRepo, time.Second*3)

	_lsHttpHandler.NewLogStreamHandler(app, logStreamUseCase)
	_lHttpHandler.NewLogHandler(app, logUseCase)
	_lsWsHandler.NewLogStreamWsHandler(app)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Send([]byte("pong"))
	})

	log.Fatal(app.Listen("127.0.0.1:7000"))
}
