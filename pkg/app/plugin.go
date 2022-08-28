package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/log-rush/distribution-server/domain"
)

type BasePlugin interface {
	Name() string
	OnInit(context *Context)
	OnAfterServe(context *Context)
	OnAfterClose(context *Context)
	OnDeInit(context *Context)
}

type Plugin interface {
	BasePlugin
	HandleLog(log domain.Log)
	SetupRouter(router fiber.Router, context *Context)
	AppendLogger(context *Context) domain.Logger
}

type HandleLog = func(log domain.Log)
type SetupRouter = func(router fiber.Router, context *Context)
type AppendLogger = func(context *Context) domain.Logger

type LogPlugin interface {
	BasePlugin
	HandleLog(log domain.Log)
}

type RouterPlugin interface {
	BasePlugin
	SetupRouter(router fiber.Router, context *Context)
}

type LoggerPlugin interface {
	BasePlugin
	AppendLogger(context *Context) domain.Logger
}
