package domain

import "time"

type Config struct {
	Timeout                  time.Duration
	LogWorkers               int
	MaxAmountOfStoredLogs    int
	ClientCheckInterval      time.Duration
	MaxClientResponseLatency time.Duration
}
