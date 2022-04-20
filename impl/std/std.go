// NOTE: This package differs from log.Std in the following ways:
//
// - Renamed Std struct to Logger.
//
// - Added Logger.logger field.
//
// - Added New() method.
//
// - Changed Prefix(), log() and logf() methods do not use the stdlog "log" package
// directly. Instead they use (l *Logger).
package std

import (
	"fmt"
	stdlog "log"
	"strings"

	"github.com/roninzo/log"
	"github.com/roninzo/log/levels"
)

// Logger is a wrapper around an instance of a logger from the Go standard
// library.
type Logger struct {
	logger *stdlog.Logger
	level  levels.Type
}

// New creates an instance of std.Logger that wraps a logger from the standard
// library. It takes a logger from the standard library as an argument. This
// is now it differs from log.NewStandard() which is used by default. The
// logger is configurable.
func New(lgr *stdlog.Logger) *Logger {
	return &Logger{
		logger: lgr,
		level:  levels.Info,
	}
}

// NewStandard sets up a basic logger using the general one provided in the Go
// standard library.
func NewStandard() *Logger {
	return &Logger{
		logger: stdlog.Default(),
		level:  levels.Info,
	}
}

func (l *Logger) Named(name string) *Logger {
	prefix := log.PrefixedWithSuffix(l.Prefix(), name)
	writer := l.logger.Writer()
	flags := l.logger.Flags()
	logger := stdlog.New(writer, prefix, flags)
	return &Logger{
		logger: logger,
		level:  l.level,
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
		l.logger.SetPrefix(prefix[0])
	}
	return l.logger.Prefix()
}

func (l *Logger) Level(level ...levels.Type) levels.Type {
	if len(level) > 0 {
		l.level = level[0]
	}
	return l.level
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
	if level < l.level { // Trace(0) < Info(2) => no logging
		return
	}
	template, args := l.unmapped(msg...)                   // fields in last position
	template, args = l.levelledf(level, template, args...) // level in first position
	template = ln(template)
	switch level {
	case levels.Panic:
		l.logger.Panicf(template, args...)
	case levels.Fatal:
		l.logger.Fatalf(template, args...)
	default:
		l.logger.Printf(template, args...)
	}
}

func (l Logger) logf(level levels.Type, template string, args ...interface{}) {
	if level < l.level { // Trace(0) < Info(2) => no logging
		return
	}
	template, args = l.unmappedf(template, args...)        // fields in last position
	template, args = l.levelledf(level, template, args...) // level in first position
	template = ln(template)
	switch level {
	case levels.Panic:
		l.logger.Panicf(template, args...)
	case levels.Fatal:
		l.logger.Fatalf(template, args...)
	default:
		l.logger.Printf(template, args...)
	}
}

func (l Logger) levelledf(level levels.Type, template string, args ...interface{}) (string, []interface{}) {
	return "%s " + template, append([]interface{}{level.Bracket()}, args...)
}

func (l Logger) unmapped(args ...interface{}) (string, []interface{}) {
	var fields log.Map
	args, fields = log.ParseArgs(args...)
	args = []interface{}{fmt.Sprint(args...)}
	template := "%s"
	if len(fields) > 0 {
		return template + " %s", append(args, l.unmap(fields))
	}
	return template, args
}

func (l Logger) unmappedf(template string, args ...interface{}) (string, []interface{}) {
	var fields log.Map
	args, fields = log.ParseArgs(args...)
	args = []interface{}{fmt.Sprintf(template, args...)}
	template = "%s"
	if len(fields) > 0 {
		return template + " %s", append(args, l.unmap(fields))
	}
	return template, args
}

func (l Logger) unmap(fields log.Map) string {
	var ret string
	for key, val := range fields {
		ret += fmt.Sprintf("[%s=%s] ", key, val)
	}
	return strings.TrimSuffix(ret, " ")
}

func ln(s string) string {
	if n := len(s); n == 0 || s[n-1] != '\n' {
		s += "\n"
	}
	return s
}
