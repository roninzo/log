package hclog

import (
	"fmt"
	"os"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/roninzo/log"
	"github.com/roninzo/log/levels"
)

// Logger is a logger that wraps the hc-log logger and has it conform to the
// log.Logger interface
type Logger struct {
	logger hclog.Logger
}

// New takes an existing hc-log logger and uses that for logging
func New(lgr hclog.Logger) *Logger {
	return &Logger{
		logger: lgr,
	}
}

// NewStandard returns a logger with a hc-log standard logger which it
// instantiates
func NewStandard() *Logger {
	return &Logger{
		logger: hclog.Default(),
	}
}

func (l *Logger) Named(name string) *Logger {
	return &Logger{
		logger: l.logger.Named(name),
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
		l.logger = l.logger.ResetNamed(prefix[0])
	}
	return l.logger.Name()
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
	if level < l.getLevel() { // Trace(0) < Info(2) => no logging
		return
	}
	var fields log.Map
	args, fields = log.ParseArgs(args...)
	msg := fmt.Sprint(args...)
	args = []interface{}{}
	if len(fields) > 0 {
		for key, val := range fields {
			args = append(args, key, val)
		}
	}
	switch level {
	case levels.Panic:
		l.logger.Error(msg, args...)
		panic(msg)
	case levels.Fatal:
		l.logger.Error(msg, args...)
		os.Exit(1)
	case levels.Error:
		l.logger.Error(msg, args...)
	case levels.Warn:
		l.logger.Warn(msg, args...)
	case levels.Info:
		l.logger.Info(msg, args...)
	case levels.Debug:
		l.logger.Debug(msg, args...)
	case levels.Trace:
		l.logger.Trace(msg, args...)
	default:
		l.logger.Info(msg, args...)
	}
}

func (l Logger) logf(level levels.Type, template string, args ...interface{}) {
	if level < l.getLevel() { // Trace(0) < Info(2) => no logging
		return
	}
	var fields log.Map
	args, fields = log.ParseArgs(args...)
	msg := fmt.Sprintf(template, args...)
	args = []interface{}{}
	if len(fields) > 0 {
		for key, val := range fields {
			args = append(args, key, val)
		}
	}
	switch level {
	case levels.Panic:
		l.logger.Error(msg, args...) // l.logger.With(unmap(fields)...).Error(msg)
		panic(msg)
	case levels.Fatal:
		l.logger.Error(msg, args...)
		os.Exit(1)
	case levels.Error:
		l.logger.Error(msg, args...)
	case levels.Warn:
		l.logger.Warn(msg, args...)
	case levels.Info:
		l.logger.Info(msg, args...)
	case levels.Debug:
		l.logger.Debug(msg, args...)
	case levels.Trace:
		l.logger.Trace(msg, args...)
	default:
		l.logger.Info(msg, args...)
	}
}

func (l Logger) getLevel() levels.Type {
	switch {
	case l.logger.IsTrace():
		return levels.Trace
	case l.logger.IsDebug():
		return levels.Debug
	case l.logger.IsInfo():
		return levels.Info
	case l.logger.IsWarn():
		return levels.Warn
	case l.logger.IsError():
		return levels.Error
	default: // case hclog.NoLevel, hclog.Off:
		return levels.Silent
	}
}

func (l *Logger) setLevel(level levels.Type) {
	switch level {
	case levels.Trace:
		l.logger.SetLevel(hclog.Trace)
	case levels.Debug:
		l.logger.SetLevel(hclog.Debug)
	case levels.Info:
		l.logger.SetLevel(hclog.Info)
	case levels.Warn:
		l.logger.SetLevel(hclog.Warn)
	case levels.Error:
		l.logger.SetLevel(hclog.Error)
	case levels.Fatal:
		l.logger.SetLevel(hclog.Error)
	case levels.Panic:
		l.logger.SetLevel(hclog.Error)
	default: // levels.Silent
		l.logger.SetLevel(hclog.Off)
	}
}
