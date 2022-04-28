package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	http_common "github.com/log-rush/simple-server/common/delivery/http"
	"github.com/log-rush/simple-server/domain"
)

type LogStreamHttpHandler struct {
	lsu domain.LogStreamUseCase
}

type RegisterRequest struct {
	Alias string `json:"alias"`
}

type UnregisterRequest struct {
	ID string `json:"id"`
}

type LogStreamsResponse struct {
	Streams []LogStreamResponse `json:"streams"`
}

type LogStreamResponse struct {
	ID    string `json:"id"`
	Alias string `json:"alias"`
}

func NewLogStreamHandler(app *fiber.App, us domain.LogStreamUseCase) {
	handler := &LogStreamHttpHandler{
		lsu: us,
	}

	app.Get("/streams", handler.ListStreams)

	streams := app.Group("/stream")
	streams.Post("/register", handler.RegisterStream)
	streams.Post("/unregister", handler.UnregisterStream)
}

// Register Streams godoc
// @ID register-stream
// @Router /stream/register [post]
// @Tags logstream
// @Description create a new logstream on the server so that client can subscribe to it
// @Summary register a logstream
// @Accept json
// @Produce json
// @Success 200 {object} LogStreamResponse
// @Failure 409 {object} http_common.ErrorResponse
// @Failure 422 {object} http_common.ErrorResponse
// @Failure 500 {object} http_common.ErrorResponse
func (h *LogStreamHttpHandler) RegisterStream(c *fiber.Ctx) error {
	ctx := c.Context()
	payload := RegisterRequest{}

	if err := c.BodyParser(&payload); err != nil {
		c.JSON(http_common.ErrorResponse{Message: err.Error()})
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	//TODO: validate

	stream, err := h.lsu.RegisterStream(ctx, payload.Alias)
	if err != nil {
		c.JSON(http_common.ErrorResponse{Message: err.Error()})
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
// @Success 200 {object} http_common.SuccessResponse
// @Failure 404 {object} http_common.ErrorResponse
// @Failure 422 {object} http_common.ErrorResponse
// @Failure 500 {object} http_common.ErrorResponse
func (h *LogStreamHttpHandler) UnregisterStream(c *fiber.Ctx) error {
	ctx := c.Context()
	payload := UnregisterRequest{}

	if err := c.BodyParser(&payload); err != nil {
		c.JSON(http_common.ErrorResponse{Message: err.Error()})
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	//TODO: validate

	err := h.lsu.UnregisterStream(ctx, payload.ID)
	if err != nil {
		c.JSON(http_common.ErrorResponse{Message: err.Error()})
		return c.SendStatus(getStatusCode(err))
	}

	c.JSON(http_common.SuccessResponse{Success: true})
	return c.SendStatus(http.StatusOK)
}

// List Streams godoc
// @ID list-stream
// @Router /streams [get]
// @Tags logstream
// @Description get all available logstreams on the server
// @Summary list all streams
// @Produce json
// @Success 200 {object} LogStreamsResponse
// @Failure 500 {object} http_common.ErrorResponse
func (h *LogStreamHttpHandler) ListStreams(c *fiber.Ctx) error {
	ctx := c.Context()

	streams, err := h.lsu.GetAvailableStreams(ctx)
	if err != nil {
		c.JSON(http_common.ErrorResponse{Message: err.Error()})
		return c.SendStatus(getStatusCode(err))
	}

	streamsDto := make([]LogStreamResponse, len(streams))
	for idx, stream := range streams {
		streamsDto[idx] = LogStreamResponse{stream.ID, stream.Alias}
	}
	c.JSON(LogStreamsResponse{streamsDto})
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
