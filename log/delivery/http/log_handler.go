package http

import (
	"net/http"

	"github.com/fabiankachlock/log-rush-simple-server/domain"
	"github.com/gofiber/fiber/v2"
)

type LogHttpHandler struct {
	useCase domain.LogUseCase
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type LogRequest struct {
	Stream    string `json:"stream"`
	Log       string `json:"log"`
	Timestamp int    `json:"timestamp"`
}

type LogBatchRequest struct {
	Stream string `json:"stream"`
	Logs   [](struct {
		Timestamp int    `json:"timestamp"`
		Log       string `json:"log"`
	}) `json:"logs"`
}

func NewLogHandler(app *fiber.App, us domain.LogUseCase) {
	handler := &LogHttpHandler{
		useCase: us,
	}

	app.Post("/log", handler.Log)
	app.Post("/batch", handler.LogBatch)
}

// Log godoc
// @ID log
// @Router /log [post]
// @Tags log
// @Description add a new log on a logstream
// @Summary push a log
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *LogHttpHandler) Log(c *fiber.Ctx) error {
	ctx := c.Context()
	payload := LogRequest{}

	if err := c.BodyParser(&payload); err != nil {
		c.JSON(ErrorResponse{err.Error()})
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	//TODO: validate

	err := h.useCase.SendLog(ctx, payload.Stream, &domain.Log{
		Message:   payload.Log,
		Stream:    payload.Stream,
		TimeStamp: payload.Timestamp,
	})
	if err != nil {
		c.JSON(ErrorResponse{err.Error()})
		return c.SendStatus(getStatusCode(err))
	}

	c.JSON(SuccessResponse{true})
	return c.SendStatus(http.StatusOK)
}

// LogBatch godoc
// @ID batch
// @Router /batch [post]
// @Tags log
// @Description add a bunch of logs at once to a stream
// @Summary push multiple logs at once
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *LogHttpHandler) LogBatch(c *fiber.Ctx) error {
	ctx := c.Context()
	payload := LogBatchRequest{}

	if err := c.BodyParser(&payload); err != nil {
		c.JSON(ErrorResponse{err.Error()})
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	//TODO: validate
	logs := make([]domain.Log, len(payload.Logs))
	for _, log := range payload.Logs {
		logs = append(logs, domain.Log{
			Message:   log.Log,
			Stream:    payload.Stream,
			TimeStamp: log.Timestamp,
		})
	}

	err := h.useCase.SendLogBatch(ctx, payload.Stream, &logs)
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
