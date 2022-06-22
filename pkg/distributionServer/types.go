package distributionServer

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/log-rush/distribution-server/domain"
)

type Log = domain.Log
type Config struct {
	Version                  string
	Name                     string
	ServerID                 string
	Timeout                  time.Duration
	LogWorkers               int
	LogsChannelBuffer        int
	MaxAmountOfStoredLogs    int
	ClientCheckInterval      time.Duration
	MaxClientResponseLatency time.Duration
	Production               bool
	Host                     string
	Port                     int
}
type LogHandler = func(streamId string, log Log)

type server struct {
	server *fiber.App
	config Config
}
