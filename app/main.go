package main

import (
	"log"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_cRepo "github.com/log-rush/simple-server/clients/repository/memory"
	_cUseCase "github.com/log-rush/simple-server/clients/usecase"
	_cfHttpHandler "github.com/log-rush/simple-server/config/delivery/http"
	_ "github.com/log-rush/simple-server/docs"
	"github.com/log-rush/simple-server/domain"
	_lHttpHandler "github.com/log-rush/simple-server/log/delivery/http"
	_lRepo "github.com/log-rush/simple-server/log/repository/memory"
	_lUseCase "github.com/log-rush/simple-server/log/usecase"
	_lsHttpHandler "github.com/log-rush/simple-server/logstream/delivery/http"
	_lsWsHandler "github.com/log-rush/simple-server/logstream/delivery/ws"
	_lsRepo "github.com/log-rush/simple-server/logstream/repository/memory"
	_lsUseCase "github.com/log-rush/simple-server/logstream/usecase"
	_sRepo "github.com/log-rush/simple-server/subscriptions/repository/memory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	mainLogger := CreateLogger()
	fiberLogger := mainLogger.Named("[server]")
	config := domain.Config{
		Name:                     "Simple log-rush distribution server",
		Version:                  "0.0.0",
		Timeout:                  time.Millisecond * 500,
		LogWorkers:               runtime.NumCPU() * 4,
		MaxAmountOfStoredLogs:    5,
		ClientCheckInterval:      time.Minute * 5,
		MaxClientResponseLatency: time.Second * 5,
	}

	app.Use(func(c *fiber.Ctx) error {
		fiberLogger.Infof("[%s] [%s] - %s", c.IP(), c.Method(), c.Path())
		err := c.Next()
		if err != nil {
			fiberLogger.Errorf("[%s] [%s] failed executing request: %s", c.IP(), c.Method(), err.Error())
		} else if c.Response().StatusCode() >= 400 {
			fiberLogger.Warnf("[%s] [%s] sending error response %d", c.IP(), c.Method(), c.Response().StatusCode())
		}
		return err
	})
	app.Use(cors.New())
	app.Use(recover.New())

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	logRepo := _lRepo.NewLogRepository(config.MaxAmountOfStoredLogs)
	logStreamRepo := _lsRepo.NewLogStreamRepository()
	clientsRepo := _cRepo.NewClientsMemoryrepository()
	subscriptionsRepo := _sRepo.NewSubscriptionsRepository(logStreamRepo)

	logStreamUseCase := _lsUseCase.NewLogStreamUseCase(logStreamRepo, subscriptionsRepo, config.LogWorkers, config.Timeout, mainLogger.Named("[logstream]"))
	logUseCase := _lUseCase.NewLogUseCase(logRepo, logStreamRepo, config.Timeout, mainLogger.Named("[logs]"))
	clientsUseCase := _cUseCase.NewClientsUseCase(clientsRepo, subscriptionsRepo, config.ClientCheckInterval, config.MaxClientResponseLatency, config.Timeout, mainLogger.Named("[clients]"))

	_lsHttpHandler.NewLogStreamHandler(app, logStreamUseCase)
	_lHttpHandler.NewLogHandler(app, logUseCase)
	_lsWsHandler.NewLogStreamWsHandler(app, clientsUseCase, mainLogger.Named("[websockets]"))
	_cfHttpHandler.NewConfigHttpHandler(app, config.Version, config.Name)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Send([]byte("pong"))
	})

	log.Fatal(app.Listen("127.0.0.1:7000"))
}

func CreateLogger() zap.SugaredLogger {
	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:   false,
		DisableCaller: true,
		Encoding:      "console",
		OutputPaths:   []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "",
			FunctionKey:    "",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
	}

	l, err := config.Build()
	defer l.Sync()
	if err != nil {
		log.Fatalln(err)
	}
	return *l.Sugar()
}
