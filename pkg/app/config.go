package app

import "time"

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
