package zap

import (
	"fmt"

	"github.com/roninzo/log"
	"github.com/roninzo/log/levels"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultAtomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
)

// New creates an instance of Zap that wraps a zap unsugared logger. It takes
// a preconfigured zap logger as an argument.
func New(lgr *zap.Logger, atoms ...zap.AtomicLevel) *Logger {
	atom := defaultAtomicLevel
	if len(atoms) > 0 {
		atom = atoms[0]
	}
	return &Logger{
		logger: lgr,
		atom:   &atom,
	}
}

// Logger is a wrapper about a Zap logger that implements the log.Logger interface.
type Logger struct {
	logger *zap.Logger
	prefix string

	// An AtomicLevel is an atomically changeable, dynamic logging level.
	// It lets you safely change the log level of a tree of loggers (the root
	// logger and any children created by adding context) at runtime.
	//
	// The AtomicLevel itself is an http.Handler that serves a JSON endpoint to
	// alter its level.
	//
	// AtomicLevels must be created with the NewAtomicLevel constructor to
	// allocate their internal atomic pointer.
	//
	// https://pkg.go.dev/go.uber.org/zap?utm_source=godoc#AtomicLevel
	atom *zap.AtomicLevel
}

func (l *Logger) Named(name string) *Logger {
	// TODO: Fix Named(using internal log.name) vs Prefix(using external prefix) inconsistencies.
	logger := *l.logger // otherwise, all zap loggers created will be one and the same.
	level := *l.atom    // otherwise, all zap loggers created will change level at the same time.
	return &Logger{
		logger: &logger, // l.logger.Named(name),
		atom:   &level,
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

func (l Logger) log(level levels.Type, args ...interface{}) {
	var fields log.Map
	args, fields = log.ParseArgs(args...)
	args = l.prefixed(args...)
	if ce := l.logger.Check(l.intLevel(level), fmt.Sprint(args...)); ce != nil {
		ce.Write(l.unmap(fields)...)
	}
}

func (l Logger) logf(level levels.Type, template string, args ...interface{}) {
	var fields log.Map
	args, fields = log.ParseArgs(args...)
	template, args = l.prefixedf(template, args...)
	if ce := l.logger.Check(l.intLevel(level), fmt.Sprintf(template, args...)); ce != nil {
		ce.Write(l.unmap(fields)...)
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

const (
	cause = "missing zap.AtomicLevel"
	fix   = "call zap.New() with *zap.AtomicLevel argument"
)

func (l Logger) getLevel() levels.Type {
	if l.atom != nil {
		switch l.atom.Level() {
		case zapcore.FatalLevel:
			return levels.Fatal
		case zapcore.PanicLevel:
			return levels.Panic
		case zapcore.DPanicLevel:
			return levels.Panic
		case zapcore.ErrorLevel:
			return levels.Error
		case zapcore.WarnLevel:
			return levels.Warn
		case zapcore.InfoLevel:
			return levels.Info
		case zapcore.DebugLevel:
			return levels.Debug
		case zapcore.DebugLevel - 1:
			return levels.Trace
		default:
			return levels.Silent
		}
	}
	// // TODO: Fix the following code and decomission fatal ending.
	// switch {
	// case l.logger.Core().Enabled(zapcore.FatalLevel):
	// 	return levels.Fatal
	// case l.logger.Core().Enabled(zapcore.PanicLevel):
	// 	return levels.Panic
	// case l.logger.Core().Enabled(zapcore.DPanicLevel):
	// 	return levels.Panic
	// case l.logger.Core().Enabled(zapcore.ErrorLevel):
	// 	return levels.Error
	// case l.logger.Core().Enabled(zapcore.WarnLevel):
	// 	return levels.Warn
	// case l.logger.Core().Enabled(zapcore.InfoLevel):
	// 	return levels.Info
	// case l.logger.Core().Enabled(zapcore.DebugLevel):
	// 	return levels.Debug
	// case l.logger.Core().Enabled(zapcore.DebugLevel - 1):
	// 	return levels.Trace
	// default:
	// 	return levels.Silent
	// }
	l.Fatal("could not get current log level", log.Map{"cause": cause, "fix": fix})
	return levels.Silent
}

func (l *Logger) setLevel(level levels.Type) {
	if l.atom != nil {
		l.atom.SetLevel(l.intLevel(level))
		return
	}
	// // TODO: Fix the following code and decomission fatal ending.
	// l.logger.Core().SetLevel(intLevel(level)):
	l.Fatalf("could not change log level: %s", level, log.Map{"cause": cause, "fix": fix})
}

func (l Logger) unmap(fields log.Map) []zapcore.Field {
	var ret []zapcore.Field
	for key, val := range fields {
		ret = append(ret, zap.Any(key, val))
	}
	return ret
}

func (l Logger) intLevel(level levels.Type) zapcore.Level {
	switch level {
	case levels.Fatal:
		return zap.FatalLevel
	case levels.Panic:
		return zap.PanicLevel // , levels.DPanic
	case levels.Error:
		return zap.ErrorLevel
	case levels.Warn:
		return zap.WarnLevel
	case levels.Info:
		return zap.InfoLevel
	case levels.Debug:
		return zap.DebugLevel
	case levels.Trace:
		return zap.DebugLevel - 1
	default:
		return zap.InfoLevel
	}
}
