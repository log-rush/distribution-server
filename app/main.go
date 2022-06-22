package main

import (
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_cRepo "github.com/log-rush/distribution-server/clients/repository/memory"
	_cUseCase "github.com/log-rush/distribution-server/clients/usecase"
	_cfHttpHandler "github.com/log-rush/distribution-server/config/delivery/http"
	_ "github.com/log-rush/distribution-server/docs"
	"github.com/log-rush/distribution-server/domain"
	_lHttpHandler "github.com/log-rush/distribution-server/log/delivery/http"
	_lRepo "github.com/log-rush/distribution-server/log/repository/memory"
	_lUseCase "github.com/log-rush/distribution-server/log/usecase"
	_lsHttpHandler "github.com/log-rush/distribution-server/logstream/delivery/http"
	_lsWsHandler "github.com/log-rush/distribution-server/logstream/delivery/ws"
	_lsRepo "github.com/log-rush/distribution-server/logstream/repository/memory"
	_lsUseCase "github.com/log-rush/distribution-server/logstream/usecase"
	_sRepo "github.com/log-rush/distribution-server/subscriptions/repository/memory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title log-rush-distribution-server
// @version 1.0
// @description this is a log-rush log distribution server
// @contact.name API Support
// @contact.email contact@fabiankachlock.dev
// @license.name MIT
// @license.url https://github.com/log-rush/distribution-server/blob/main/LICENSE
// @host localhost:7000
// @BasePath /
// @Tag.name logstream
// @Tag.description all endpoints for logstreams
// @Tag.name log
// @Tag.description all endpoints for logs
func main() {
	isProd := flag.Bool("prod", false, "Indicates if the server should be run in production mode")
	app := fiber.New(fiber.Config{
		Prefork:           *isProd,
		EnablePrintRoutes: true,
	})

	mainLogger := CreateLogger()
	fiberLogger := mainLogger.Named("[server]")
	config := domain.Config{
		Name:                     "Simple log-rush distribution server",
		Version:                  "0.0.0",
		ServerID:                 "dev-server",
		Timeout:                  time.Millisecond * 500,
		LogWorkers:               runtime.NumCPU() * 4,
		MaxAmountOfStoredLogs:    100,
		LogsChannelBuffer:        16,
		ClientCheckInterval:      time.Minute * 1,
		MaxClientResponseLatency: time.Second * 2,
	}

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

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

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/test", func(c *fiber.Ctx) error {
		panic("test")
	})

	logRepo := _lRepo.NewLogRepository(config.MaxAmountOfStoredLogs)
	logStreamRepo := _lsRepo.NewLogStreamRepository(config.LogsChannelBuffer)
	clientsRepo := _cRepo.NewClientsMemoryrepository()
	subscriptionsRepo := _sRepo.NewSubscriptionsRepository(logStreamRepo)

	logStreamUseCase := _lsUseCase.NewLogStreamUseCase(logStreamRepo, subscriptionsRepo, config.LogWorkers, config.Timeout, mainLogger.Named("[logstream]"))
	logUseCase := _lUseCase.NewLogUseCase(logRepo, logStreamRepo, config.Timeout, mainLogger.Named("[logs]"))
	clientsUseCase := _cUseCase.NewClientsUseCase(clientsRepo, subscriptionsRepo, logRepo, config.ClientCheckInterval, config.MaxClientResponseLatency, config.Timeout, mainLogger.Named("[clients]"))

	_lsHttpHandler.NewLogStreamHandler(app, logStreamUseCase)
	_lHttpHandler.NewLogHandler(app, logUseCase)
	_lsWsHandler.NewLogStreamWsHandler(app, clientsUseCase, mainLogger.Named("[websockets]"))
	_cfHttpHandler.NewConfigHttpHandler(app, config.Version, config.Name, config.ServerID)

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
