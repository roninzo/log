package logger

import (
	"time"
)

type Config struct {
	IgnoreRecordNotFoundError bool          // Ignore ErrRecordNotFound Error
	SlowThreshold             time.Duration // Threshold for reporting slow SQL queries
}

var defaultConfig = Config{
	IgnoreRecordNotFoundError: false,
	SlowThreshold:             200 * time.Millisecond,
}
