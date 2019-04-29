package types

import "time"

// A Timeframe denotes a start and end point in time.
type Timeframe struct {
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}
