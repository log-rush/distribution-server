package devkit

import (
	"github.com/gofiber/fiber/v2"
	"github.com/log-rush/distribution-server/domain"
	"github.com/log-rush/distribution-server/pkg/app"
)

type Plugin struct {
	name          string
	LogHandler    app.HandleLog
	RouterHandler app.SetupRouter
	LoggerHandler app.AppendLogger
}

func (p Plugin) HandleLog(log domain.Log) {
	p.LogHandler(log)
}

func (p Plugin) SetupRouter(router fiber.Router, context *app.Context) {
	p.RouterHandler(router, context)
}

func (p Plugin) AppendLogger(context *app.Context) domain.Logger {
	return p.LoggerHandler(context)
}

func (p Plugin) Name() string {
	return p.name
}

func NewPlugin(
	name string,
	logHandler app.HandleLog,
	routerHandler app.SetupRouter,
	appendLogger app.AppendLogger,
) *Plugin {
	p := Plugin{
		name:          name,
		LogHandler:    logHandler,
		RouterHandler: routerHandler,
		LoggerHandler: appendLogger,
	}

	return &p
}

func NewLogPlugin(
	name string,
	logHandler app.HandleLog,
) *Plugin {
	p := Plugin{
		name:          name,
		LogHandler:    logHandler,
		RouterHandler: nil,
		LoggerHandler: nil,
	}

	return &p
}

func NewRouterPlugin(
	name string,
	routerHandler app.SetupRouter,
) *Plugin {
	p := Plugin{
		name:          name,
		RouterHandler: routerHandler,
		LogHandler:    nil,
		LoggerHandler: nil,
	}

	return &p
}

func NewLoggerPlugin(
	name string,
	appendLogger app.AppendLogger,
) *Plugin {
	p := Plugin{
		name:          name,
		RouterHandler: nil,
		LogHandler:    nil,
		LoggerHandler: appendLogger,
	}

	return &p
}
