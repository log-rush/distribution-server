package distributionServer

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_cRepo "github.com/log-rush/distribution-server/clients/repository/memory"
	_cUseCase "github.com/log-rush/distribution-server/clients/usecase"
	_cfHttpHandler "github.com/log-rush/distribution-server/config/delivery/http"
	_ "github.com/log-rush/distribution-server/docs"
	_lHttpHandler "github.com/log-rush/distribution-server/log/delivery/http"
	_lRepo "github.com/log-rush/distribution-server/log/repository/memory"
	_lUseCase "github.com/log-rush/distribution-server/log/usecase"
	_lsHttpHandler "github.com/log-rush/distribution-server/logstream/delivery/http"
	_lsWsHandler "github.com/log-rush/distribution-server/logstream/delivery/ws"
	_lsRepo "github.com/log-rush/distribution-server/logstream/repository/memory"
	_lsUseCase "github.com/log-rush/distribution-server/logstream/usecase"
	_app "github.com/log-rush/distribution-server/pkg/app"
	"github.com/log-rush/distribution-server/pkg/devkit"
	_sRepo "github.com/log-rush/distribution-server/subscriptions/repository/memory"
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
func NewServer(config _app.Config) *server {
	app := fiber.New(fiber.Config{
		Prefork:           config.Production,
		EnablePrintRoutes: !config.Production,
	})

	var mainLogger = NewMultiLogger(createLogger())
	server := server{
		server:  app,
		config:  config,
		context: _app.NewAppContext(config, app, mainLogger),
	}

	return &server
}

func (s *server) Start() {
	appContext := s.context
	app := s.server

	// setup loggers
	multiLogger := (*appContext.Logger).(*MultiLogger)
	for _, p := range *appContext.Plugins.LoggerPlugins {
		logger := p.AppendLogger(appContext)
		*multiLogger.Loggers = append(*multiLogger.Loggers, logger)
	}

	fiberLogger := (*appContext.Logger).Named("[server]")
	pluginLogger := (*appContext.Logger).Named("[plugin]")

	for _, p := range *appContext.Plugins.LoggerPlugins {
		pluginLogger.Debugf("active logger: %s", p.Name())
	}

	// recover from errors in handlers
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// log fiber requests
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

	// use cors
	app.Use(cors.New())
	// use swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// init repositories
	appContext.Repos.Log = _lRepo.NewLogRepository(appContext)
	appContext.Repos.LogStream = _lsRepo.NewLogStreamRepository(appContext)
	appContext.Repos.Clients = _cRepo.NewClientsMemoryRepository()
	appContext.Repos.Subscriptions = _sRepo.NewSubscriptionsRepository(appContext)

	// init use cases
	appContext.UseCases.LogStream = _lsUseCase.NewLogStreamUseCase(appContext)
	appContext.UseCases.Log = _lUseCase.NewLogUseCase(appContext)
	appContext.UseCases.Clients = _cUseCase.NewClientsUseCase(appContext)

	// init http handlers
	_lsHttpHandler.NewLogStreamHandler(appContext)
	_lHttpHandler.NewLogHandler(appContext)
	_lsWsHandler.NewLogStreamWsHandler(appContext)
	_cfHttpHandler.NewConfigHttpHandler(appContext)

	// setup ping handler
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Send([]byte("pong"))
	})

	// setup router plugins
	for _, plugin := range *s.context.Plugins.RouterPlugins {
		router := s.server.Group("/plugins/" + plugin.Name())
		pluginLogger.Debugf("setting up router: %s", plugin.Name())
		plugin.SetupRouter(router, s.context)
	}

	// call plugin hooks
	multiLogger.Info("calling hook: init")
	for _, p := range *appContext.Plugins.Plugins {
		p.OnInit(appContext)
	}

	defer func() {
		multiLogger.Info("calling hook: deinit")
		for _, p := range *appContext.Plugins.Plugins {
			p.OnDeInit(appContext)
		}
	}()

	app.Hooks().OnListen(func() error {
		multiLogger.Info("calling hook: afterServe")
		for _, p := range *appContext.Plugins.Plugins {
			p.OnAfterServe(appContext)
		}
		return nil
	})

	app.Hooks().OnShutdown(func() error {
		multiLogger.Info("calling hook: afterClose")
		for _, p := range *appContext.Plugins.Plugins {
			p.OnAfterClose(appContext)
		}
		return nil
	})

	// start server
	log.Fatal(s.server.Listen(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)))
}

func (s *server) Stop() error {
	return s.server.Shutdown()
}

func (s *server) UsePlugin(plugin _app.Plugin) {
	fmt.Println("using", plugin.Name())
	p := plugin.(*devkit.Plugin)
	*s.context.Plugins.Plugins = append(*s.context.Plugins.Plugins, plugin)

	if p.LogHandler != nil {
		*s.context.Plugins.LogPlugins = append(*s.context.Plugins.LogPlugins, plugin)
	}
	if p.LoggerHandler != nil {
		*s.context.Plugins.LoggerPlugins = append(*s.context.Plugins.LoggerPlugins, plugin)
	}
	if p.RouterHandler != nil {
		*s.context.Plugins.RouterPlugins = append(*s.context.Plugins.RouterPlugins, plugin)
	}
}
