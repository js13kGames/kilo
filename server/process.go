package server

import "time"

// A Process is a set of operations that can be started and stopped.
type Process interface {
	// Start starts the Processes operations.
	// Start is expected to be blocking.
	Start()

	// Stop stops the Processes operations. If deadline is non nil the implementation is expected
	// to stop gracefully no later than after the deadline expires.
	// Stop is expected to be blocking.
	Stop(deadline *time.Time)
}
