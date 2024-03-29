package main

import (
	"log"
	"runtime"
	"time"

	"github.com/log-rush/distribution-server/pkg/app"
	"github.com/log-rush/distribution-server/pkg/distributionServer"
	storageAdapterFs "github.com/log-rush/persistency-adapter-fs"
	pluginPersistency "github.com/log-rush/plugin-persistency"
	pluginServerLogs "github.com/log-rush/plugin-server-logs"
)

func main() {
	config := app.Config{
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
		BasePath:                "./_logs",
		OpenHandleTimeout:       time.Minute * 10,
		ForceUpdateOnMidnight:   false,
		DateFormat:              "02_01_06",
		GroupStreamsIntoFolders: true,
		FilePermissions:         0744,
	})
	if err != nil {
		log.Fatalf("cant init fs storage adapter\n")
	}

	plugin := pluginPersistency.NewPersistencyPlugin(pluginPersistency.Config{
		Adapter:          adapter,
		LogDelimiter:     "\n",
		StreamsBlacklist: []string{},
		StreamsWhitelist: []string{},
	}).Plugin

	server.UsePlugin(plugin)

	loggingPlugin := pluginServerLogs.NewServerLogsPlugin(pluginServerLogs.Config{
		StreamName: "server",
		Id:         "server",
		Key:        "super-secret",
	}).Plugin

	server.UsePlugin(loggingPlugin)

	// server.UsePlugin(devkit.NewPlugin("test-logger", nil, nil, func(context *app.Context) domain.Logger {
	// 	return devkit.NewLogger(func(level devkit.LogLevel, log string, args ...interface{}) {
	// 		fmt.Printf("> testlogger > [%s] %s \n", level, fmt.Sprintf(log, args...))
	// 	})
	// }))

	// server.UsePlugin(devkit.NewPlugin("test", func(log logRush.Log) {}, func(router fiber.Router) {
	// 	router.Get("/pong", func(c *fiber.Ctx) error {
	// 		return c.SendString("ping")
	// 	})
	// }))

	server.Start()
}
