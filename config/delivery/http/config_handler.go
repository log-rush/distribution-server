package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/log-rush/distribution-server/pkg/app"
	logRush "github.com/log-rush/server-devkit/v2"
)

type configHandler struct {
	version       string
	name          string
	id            string
	routerPlugins *[]logRush.Plugin
}

type InfoResponse struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Name    string `json:"name"`
}

type PluginsResponse struct {
	RouterPlugins []string `json:"routerPlugins"`
}

func NewConfigHttpHandler(context *app.Context) {
	handler := configHandler{
		version:       context.Config.Version,
		name:          context.Config.Name,
		id:            context.Config.ServerID,
		routerPlugins: context.Plugins.RouterPlugins,
	}

	context.Server.Get("/info", handler.getInfo)
	context.Server.Get("/meta/plugins", handler.getPlugins)
}

// get Info godoc
// @ID info
// @Router /info [get]
// @Tags system
// @Description get all server information
// @Summary get the server info
// @Produce json
// @Success 200 {object} InfoResponse
func (h *configHandler) getInfo(c *fiber.Ctx) error {
	return c.Status(200).JSON(InfoResponse{
		Version: h.version,
		Name:    h.name,
		ID:      h.id,
	})
}

// get Plugins godoc
// @ID meta-plugins
// @Router /meta/plugins [get]
// @Tags system
// @Description get info about the used plugins
// @Summary get activates plugins
// @Produce json
// @Success 200 {object} PluginsResponse
func (h *configHandler) getPlugins(c *fiber.Ctx) error {
	plugins := []string{}
	for _, plugin := range *h.routerPlugins {
		plugins = append(plugins, plugin.Name())
	}

	return c.Status(200).JSON(PluginsResponse{
		RouterPlugins: plugins,
	})
}
