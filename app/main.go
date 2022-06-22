package main

import (
	"runtime"
	"time"

	"github.com/log-rush/distribution-server/pkg/distributionServer"
)

func main() {
	config := distributionServer.Config{
		Name:                     "Simple log-rush distribution server",
		Version:                  "0.0.0",
		ServerID:                 "dev-server",
		Timeout:                  time.Millisecond * 500,
		LogWorkers:               runtime.NumCPU() * 4,
		MaxAmountOfStoredLogs:    100,
		LogsChannelBuffer:        16,
		ClientCheckInterval:      time.Minute * 1,
		MaxClientResponseLatency: time.Second * 2,
		Production:               false,
		Host:                     "127.0.0.1",
		Port:                     7000,
	}

	server := distributionServer.NewServer(config)
	server.Start()
}
