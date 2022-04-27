package domain

import "time"

type Config struct {
	Timeout    time.Duration
	LogWorkers int
}
