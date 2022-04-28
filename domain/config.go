package domain

import "time"

type Config struct {
	Version                  string
	Name                     string
	Timeout                  time.Duration
	LogWorkers               int
	MaxAmountOfStoredLogs    int
	ClientCheckInterval      time.Duration
	MaxClientResponseLatency time.Duration
}
