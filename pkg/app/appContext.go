package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/log-rush/distribution-server/domain"
	logRush "github.com/log-rush/server-devkit/v2"
)

type (
	AppPlugins struct {
		LogPlugins    *[]logRush.Plugin
		RouterPlugins *[]logRush.Plugin
	}

	AppRepos struct {
		Log           domain.LogRepository
		LogStream     domain.LogStreamRepository
		Clients       domain.ClientsRepository
		Subscriptions domain.SubscriptionsRepository
	}

	AppUseCases struct {
		Log       domain.LogUseCase
		LogStream domain.LogStreamUseCase
		Clients   domain.ClientsUseCase
	}

	Context struct {
		Config   Config
		Server   *fiber.App
		Logger   *domain.Logger
		Plugins  *AppPlugins
		UseCases *AppUseCases
		Repos    *AppRepos
	}
)

func NewAppContext(config Config, server *fiber.App, logger domain.Logger) *Context {
	return &Context{
		Config: config,
		Server: server,
		Logger: &logger,
		Plugins: &AppPlugins{
			LogPlugins:    &[]logRush.Plugin{},
			RouterPlugins: &[]logRush.Plugin{},
		},
		Repos: &AppRepos{
			Log:           nil,
			LogStream:     nil,
			Clients:       nil,
			Subscriptions: nil,
		},
		UseCases: &AppUseCases{
			Log:       nil,
			LogStream: nil,
			Clients:   nil,
		},
	}
}

func (a Context) RequireModule(name string) {
	panic("cant initialize because a module named" + name + "is nil")
}
