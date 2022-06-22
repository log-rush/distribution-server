package commons

import (
	"runtime/debug"

	"github.com/log-rush/distribution-server/domain"
)

func RecoverRoutine(logger *domain.Logger) {
	r := recover()
	if r != nil {
		(*logger).Errorf("error occurred in goroutine: %s", (r).(error).Error())
		debug.PrintStack()
	}
}
