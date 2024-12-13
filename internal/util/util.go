package util

import (
	log "github.com/sirupsen/logrus"
)

// LogDeferredError wraps a function that returns an error (e.g., f.Close()) and logs any error.
func LogDeferredError(fn func() error, logger *log.Logger) func() {
	return func() {
		if err := fn(); err != nil {
			if logger != nil {
				logger.Errorf("Error in deferred function: %v", err)
			} else {
				log.Errorf("Error in deferred function: %v", err)
			}
		}
	}
}
