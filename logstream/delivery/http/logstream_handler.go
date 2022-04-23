package http

import (
	"net/http"

	"github.com/fabiankachlock/log-rush-simple-server/domain"
	"github.com/gofiber/fiber/v2"
)

type LogStreamHttpHandler struct {
	useCase domain.LogStreamUseCase
}

type RegisterRequest struct {
	Alias string `json:"alias"`
}

type UnregisterRequest struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

func NewLogStreamHandler(app *fiber.App, us domain.LogStreamUseCase) {
	handler := &LogStreamHttpHandler{
		useCase: us,
	}
	streams := app.Group("/stream")

	streams.Post("/register", handler.RegisterStream)
}

// Register Streams godoc
// @ID register-stream
// @Router /stream/register [post]
// @Tags logstream
// @Description create a new logstream on the server so that client can subscribe to it
// @Summary register a logstream
// @Accept json
// @Produce json
// @Success 200 {object} domain.LogStream
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *LogStreamHttpHandler) RegisterStream(c *fiber.Ctx) error {
	ctx := c.Context()
	payload := RegisterRequest{}

	if err := c.BodyParser(&payload); err != nil {
		c.JSON(ErrorResponse{err.Error()})
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	//TODO: validate

	stream, err := h.useCase.RegisterStream(ctx, payload.Alias)
	if err != nil {
		c.JSON(ErrorResponse{err.Error()})
		return c.SendStatus(getStatusCode(err))
	}

	c.JSON(stream)
	return c.SendStatus(http.StatusOK)
}

// Unregister Streams godoc
// @ID unregister-stream
// @Router /stream/unregister [post]
// @Tags logstream
// @Description delete a logstream an close all open connection to it
// @Summary unregister a logstream
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *LogStreamHttpHandler) UnregisterStream(c *fiber.Ctx) error {
	ctx := c.Context()
	payload := UnregisterRequest{}

	if err := c.BodyParser(&payload); err != nil {
		c.JSON(ErrorResponse{err.Error()})
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	//TODO: validate

	err := h.useCase.UnregisterStream(ctx, payload.ID)
	if err != nil {
		c.JSON(ErrorResponse{err.Error()})
		return c.SendStatus(getStatusCode(err))
	}

	c.JSON(SuccessResponse{true})
	return c.SendStatus(http.StatusOK)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case domain.ErrStreamAlreadyExists:
		return http.StatusConflict
	case domain.ErrStreamNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
