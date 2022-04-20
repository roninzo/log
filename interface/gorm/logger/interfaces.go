package logger

import (
	"context"
	"time"

	"github.com/roninzo/log"
	"gorm.io/gorm/logger"
)

// Writer log writer interface
type Writer interface {
	Info(msg string, fields log.Map)
	Warn(msg string, fields log.Map)
	Error(msg string, fields log.Map)
}

// Interface logger interface
type Interface interface {
	LogMode(level logger.LogLevel) Interface
	Info(ctx context.Context, format string, args ...interface{})
	Warn(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, format string, args ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}
