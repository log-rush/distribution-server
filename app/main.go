package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/log-rush/distribution-server/domain"
	"github.com/log-rush/distribution-server/pkg/distributionServer"
)

type testPlugin struct{}

func (x testPlugin) HandleLog(stream string, log domain.Log) {
	fmt.Printf("[custom plugin]: %s | %s\n", stream, log.Message)
}

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

	x := testPlugin{}
	server.UseLogPlugin(x)

	server.Start()
}
