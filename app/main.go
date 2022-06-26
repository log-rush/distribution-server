package main

import (
	"log"
	"runtime"
	"time"

	"github.com/log-rush/distribution-server/pkg/distributionServer"
	storageAdapterFs "github.com/log-rush/persistency-adapter-fs"
	pluginPersistency "github.com/log-rush/plugin-persistency"
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

	adapter, err := storageAdapterFs.NewFSStorageAdapter(storageAdapterFs.Config{
		BasePath:              "./_logs",
		OpenHandleTimeout:     time.Minute * 10,
		ForceUpdateOnMidnight: false,
		DateFormat:            "02_01_06",
	})
	if err != nil {
		log.Fatalf("cant init fs storage adapter\n")
	}

	plugin := pluginPersistency.NewPersistencyPlugin(pluginPersistency.Config{
		Adapter:          adapter,
		LogDelimiter:     "\n",
		StreamsBlacklist: []string{},
		StreamsWhitelist: []string{},
	})

	server.UsePlugin(*plugin.Plugin)

	server.Start()
}
