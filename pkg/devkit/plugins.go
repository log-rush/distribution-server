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
	Hooks         PluginHooks
}

type PluginHooks struct {
	OnInit       func(context *app.Context)
	OnAfterServe func(context *app.Context)
	OnAfterClose func(context *app.Context)
	OnDeInit     func(context *app.Context)
}

type PluginHandlers struct {
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

func (p Plugin) OnInit(context *app.Context) {
	if p.Hooks.OnInit != nil {
		p.Hooks.OnInit(context)
	}
}

func (p Plugin) OnAfterServe(context *app.Context) {
	if p.Hooks.OnAfterServe != nil {
		p.Hooks.OnAfterServe(context)
	}
}

func (p Plugin) OnAfterClose(context *app.Context) {
	if p.Hooks.OnAfterClose != nil {
		p.Hooks.OnAfterClose(context)
	}
}

func (p Plugin) OnDeInit(context *app.Context) {
	if p.Hooks.OnDeInit != nil {
		p.Hooks.OnDeInit(context)
	}
}

func NewPlugin(
	name string,
	handlers PluginHandlers,
	hooks PluginHooks,
) *Plugin {
	p := Plugin{
		name:          name,
		LogHandler:    handlers.LogHandler,
		RouterHandler: handlers.RouterHandler,
		LoggerHandler: handlers.LoggerHandler,
		Hooks:         hooks,
	}

	return &p
}

func NewLogPlugin(
	name string,
	logHandler app.HandleLog,
	hooks PluginHooks,
) *Plugin {
	p := Plugin{
		name:          name,
		LogHandler:    logHandler,
		RouterHandler: nil,
		LoggerHandler: nil,
		Hooks:         hooks,
	}

	return &p
}

func NewRouterPlugin(
	name string,
	routerHandler app.SetupRouter,
	hooks PluginHooks,
) *Plugin {
	p := Plugin{
		name:          name,
		RouterHandler: routerHandler,
		LogHandler:    nil,
		LoggerHandler: nil,
		Hooks:         hooks,
	}

	return &p
}

func NewLoggerPlugin(
	name string,
	appendLogger app.AppendLogger,
	hooks PluginHooks,
) *Plugin {
	p := Plugin{
		name:          name,
		RouterHandler: nil,
		LogHandler:    nil,
		LoggerHandler: appendLogger,
		Hooks:         hooks,
	}

	return &p
}
