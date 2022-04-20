package logrus

import (
	"github.com/roninzo/log"
	"github.com/roninzo/log/levels"
	"github.com/sirupsen/logrus"
)

// Logger is a logger that wraps the logrus logger and has it conform to the
// log.Logger interface
type Logger struct {
	logger *logrus.Logger
	prefix string
}

// New takes an existing logrus logger and uses that for logging
func New(lgr *logrus.Logger) *Logger {
	return &Logger{
		logger: lgr,
	}
}

// NewStandard returns a logger with a logrus standard logger which it
// instantiates
func NewStandard() *Logger {
	return &Logger{
		logger: logrus.StandardLogger(),
	}
}

func (l *Logger) Named(name string) *Logger {
	logger := *l.logger // otherwise, all zap loggers created will be one and the same.
	return &Logger{
		logger: &logger,
		prefix: log.Prefixed(l.prefix, name),
	}
}

func (l *Logger) Options(funcs ...func(*Logger) *Logger) *Logger {
	for _, f := range funcs {
		f(l)
	}
	return l
}

func (l *Logger) WithLevel(level levels.Type) *Logger {
	l.Level(level)
	return l
}

func (l *Logger) WithLevelFromDebug(debug bool) *Logger {
	switch debug {
	case true:
		l.Level(levels.Debug)
	default:
		l.Level(levels.Info)
	}
	return l
}

func (l *Logger) Prefix(prefix ...string) string {
	if len(prefix) > 0 {
		l.prefix = prefix[0]
	}
	return l.prefix
}

func (l *Logger) Level(level ...levels.Type) levels.Type {
	if len(level) > 0 {
		l.setLevel(level[0])
	}
	return l.getLevel()
}

func (l Logger) Trace(msg ...interface{}) { l.log(levels.Trace, msg...) }
func (l Logger) Debug(msg ...interface{}) { l.log(levels.Debug, msg...) }
func (l Logger) Info(msg ...interface{})  { l.log(levels.Info, msg...) }
func (l Logger) Warn(msg ...interface{})  { l.log(levels.Warn, msg...) }
func (l Logger) Error(msg ...interface{}) { l.log(levels.Error, msg...) }
func (l Logger) Panic(msg ...interface{}) { l.log(levels.Panic, msg...) }
func (l Logger) Fatal(msg ...interface{}) { l.log(levels.Fatal, msg...) }

func (l Logger) Tracef(template string, args ...interface{}) { l.logf(levels.Trace, template, args...) }
func (l Logger) Debugf(template string, args ...interface{}) { l.logf(levels.Debug, template, args...) }
func (l Logger) Infof(template string, args ...interface{})  { l.logf(levels.Info, template, args...) }
func (l Logger) Warnf(template string, args ...interface{})  { l.logf(levels.Warn, template, args...) }
func (l Logger) Errorf(template string, args ...interface{}) { l.logf(levels.Error, template, args...) }
func (l Logger) Panicf(template string, args ...interface{}) { l.logf(levels.Panic, template, args...) }
func (l Logger) Fatalf(template string, args ...interface{}) { l.logf(levels.Fatal, template, args...) }

func (l Logger) log(level levels.Type, msg ...interface{}) {
	if level < l.getLevel() { // Trace(0) < Info(2) => no logging
		return
	}
	msg = l.prefixed(msg...)
	if args, fields := log.ParseArgs(msg...); len(fields) > 0 {
		switch level {
		case levels.Panic:
			l.logger.WithFields(logrus.Fields(fields)).Panic(args...)
		case levels.Fatal:
			l.logger.WithFields(logrus.Fields(fields)).Fatal(args...)
		case levels.Error:
			l.logger.WithFields(logrus.Fields(fields)).Error(args...)
		case levels.Warn:
			l.logger.WithFields(logrus.Fields(fields)).Warn(args...)
		case levels.Info:
			l.logger.WithFields(logrus.Fields(fields)).Info(args...)
		case levels.Debug:
			l.logger.WithFields(logrus.Fields(fields)).Debug(args...)
		case levels.Trace:
			l.logger.WithFields(logrus.Fields(fields)).Trace(args...)
		default:
			l.logger.WithFields(logrus.Fields(fields)).Info(args...)
		}
		return
	}
	switch level {
	case levels.Panic:
		l.logger.Panic(msg...)
	case levels.Fatal:
		l.logger.Fatal(msg...)
	case levels.Error:
		l.logger.Error(msg...)
	case levels.Warn:
		l.logger.Warn(msg...)
	case levels.Info:
		l.logger.Info(msg...)
	case levels.Debug:
		l.logger.Debug(msg...)
	case levels.Trace:
		l.logger.Trace(msg...)
	default:
		l.logger.Info(msg...)
	}
}

func (l Logger) logf(level levels.Type, template string, args ...interface{}) {
	if level < l.getLevel() { // Trace(0) < Info(2) => no logging
		return
	}
	template, args = l.prefixedf(template, args...)
	if msg, fields := log.ParseArgs(args...); len(fields) > 0 {
		switch level {
		case levels.Panic:
			l.logger.WithFields(logrus.Fields(fields)).Panicf(template, msg...)
		case levels.Fatal:
			l.logger.WithFields(logrus.Fields(fields)).Fatalf(template, msg...)
		case levels.Error:
			l.logger.WithFields(logrus.Fields(fields)).Errorf(template, msg...)
		case levels.Warn:
			l.logger.WithFields(logrus.Fields(fields)).Warnf(template, msg...)
		case levels.Info:
			l.logger.WithFields(logrus.Fields(fields)).Infof(template, msg...)
		case levels.Debug:
			l.logger.WithFields(logrus.Fields(fields)).Debugf(template, msg...)
		case levels.Trace:
			l.logger.WithFields(logrus.Fields(fields)).Tracef(template, msg...)
		default:
			l.logger.WithFields(logrus.Fields(fields)).Infof(template, msg...)
		}
		return
	}
	switch level {
	case levels.Panic:
		l.logger.Panicf(template, args...)
	case levels.Fatal:
		l.logger.Fatalf(template, args...)
	case levels.Error:
		l.logger.Errorf(template, args...)
	case levels.Warn:
		l.logger.Warnf(template, args...)
	case levels.Info:
		l.logger.Infof(template, args...)
	case levels.Debug:
		l.logger.Debugf(template, args...)
	case levels.Trace:
		l.logger.Tracef(template, args...)
	default:
		l.logger.Infof(template, args...)
	}
}

func (l Logger) prefixed(msg ...interface{}) []interface{} {
	if l.prefix != "" {
		return append([]interface{}{l.prefix + ": "}, msg...)
	}
	return msg
}

func (l Logger) prefixedf(template string, args ...interface{}) (string, []interface{}) {
	if l.prefix != "" {
		return "%s: " + template, append([]interface{}{l.prefix}, args...)
	}
	return template, args
}

func (l Logger) getLevel() levels.Type {
	switch l.logger.GetLevel() {
	case logrus.PanicLevel:
		return levels.Panic
	case logrus.FatalLevel:
		return levels.Fatal
	case logrus.ErrorLevel:
		return levels.Error
	case logrus.WarnLevel:
		return levels.Warn
	case logrus.InfoLevel:
		return levels.Info
	case logrus.DebugLevel:
		return levels.Debug
	case logrus.TraceLevel:
		return levels.Trace
	default:
		return levels.Info
	}
}

func (l *Logger) setLevel(level levels.Type) {
	switch level {
	case levels.Panic:
		l.logger.SetLevel(logrus.PanicLevel)
	case levels.Fatal:
		l.logger.SetLevel(logrus.FatalLevel)
	case levels.Error:
		l.logger.SetLevel(logrus.ErrorLevel)
	case levels.Warn:
		l.logger.SetLevel(logrus.WarnLevel)
	case levels.Info:
		l.logger.SetLevel(logrus.InfoLevel)
	case levels.Debug:
		l.logger.SetLevel(logrus.DebugLevel)
	case levels.Trace:
		l.logger.SetLevel(logrus.TraceLevel)
	default:
		l.logger.SetLevel(logrus.InfoLevel)
	}
}
