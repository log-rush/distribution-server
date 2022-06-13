package commons

import (
	"runtime/debug"

	"github.com/log-rush/simple-server/domain"
)

func RecoverRoutine(logger *domain.Logger) {
	r := recover()
	if err := r.(error); err != nil {
		(*logger).Errorf("error occurred in goroutine: %s", err.Error())
		debug.PrintStack()
	}
}
