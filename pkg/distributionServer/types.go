package distributionServer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/log-rush/distribution-server/pkg/app"
)

type server struct {
	server  *fiber.App
	config  app.Config
	context *app.Context
}
