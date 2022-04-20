package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/roninzo/log"
	"github.com/roninzo/log/levels"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// Gorm logger via composition.
type Logger struct {
	logger log.Logger // Underlying Logger instance.
	Config            // Extra Gorm config.
}

// Constructor.
func New(lgr log.Logger, configs ...Config) Interface {
	config := defaultConfig
	if len(configs) > 0 {
		config = configs[0]
	}
	return &Logger{lgr, config}
}

// Info print info.
func (l Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, data...), log.Map{"source": utils.FileWithLineNum()})
}

// Warn print warn messages.
func (l Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warn(fmt.Sprintf(msg, data...), log.Map{"source": utils.FileWithLineNum()})
}

// Error print error messages.
func (l Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, data...), log.Map{"source": utils.FileWithLineNum()})
}

// Trace print sql message.
func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logger.Level() >= levels.Silent { // Skip Silent
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.logger.Level() <= levels.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError): // Error/Warn/Info
		sql, rows := fc()
		l.logger.Error(traceErrMsg(utils.FileWithLineNum(), sql, rows, elapsed, err))
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.logger.Level() <= levels.Warn: // Warn/Info
		sql, rows := fc()
		l.logger.Warn(traceWarnMsg(utils.FileWithLineNum(), sql, rows, elapsed, l.SlowThreshold))
	case l.logger.Level() == levels.Info: // Info
		sql, rows := fc()
		l.logger.Info(traceMsg(utils.FileWithLineNum(), sql, rows, elapsed))
	}
}
