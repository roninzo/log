package logger

import (
	"github.com/roninzo/log/levels"
	"gorm.io/gorm/logger"
)

// LogMode log mode
func (l *Logger) LogMode(level logger.LogLevel) Interface {
	derefdLogger := *l
	newlogger := &derefdLogger
	newlogger.logger.Level(toInternalLevel(level))
	return newlogger
}

// Converts logger.LogLevel to levels.Type.
func toInternalLevel(level logger.LogLevel) levels.Type {
	switch level {
	case logger.Silent:
		return levels.Silent
	case logger.Error:
		return levels.Error
	case logger.Warn:
		return levels.Warn
	case logger.Info:
		return levels.Info
	default:
		return levels.Info
	}
}
