package http

import "github.com/gofiber/fiber/v2"

type configHandler struct {
	version string
	name    string
	id      string
}

type InfoResponse struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Name    string `json:"name"`
}

func NewConfigHttpHandler(app *fiber.App, version, name, id string) {
	handler := configHandler{version, name, id}

	app.Get("/info", handler.getInfo)
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
