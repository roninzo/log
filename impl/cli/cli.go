// This package provides a reference implementation for a CLI logger, where the
// output is written to the console. If you need a CLI logger that is fairly
// different from this one please feel free to create another CLI implementation
// and you can fork this one as a starting point. Not everyone needs to use this
// logger CLI implementation and it does not need to have all features.

// TODO: Add i18n support for level labels
// TODO: Add a mutex to lock writing output

package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/roninzo/log"
	"github.com/roninzo/log/levels"
)

// Logger provides a CLI based logger. Log messages are written to the CLI as
// terminal style output.
type Logger struct {
	// Sets the current logging level
	level levels.Type

	// Sets the current logging prefix
	prefix string

	// Outputs used for each of the levels. Provides a writer
	// to write messages to io.Writer.
	TraceOutput io.Writer
	DebugOutput io.Writer
	InfoOutput  io.Writer
	WarnOutput  io.Writer
	ErrorOutput io.Writer
	PanicOutput io.Writer
	FatalOutput io.Writer

	// Colors used for each of the levels. Note, the Info level intentionally
	// does not have a color. It will use the tty default.
	TraceColor *color.Color
	DebugColor *color.Color
	WarnColor  *color.Color
	ErrorColor *color.Color
	PanicColor *color.Color
	FatalColor *color.Color
}

// NewStandard creates a default CLI logger
func NewStandard() *Logger {
	return &Logger{
		// Note, stderr is used for all non-info messages by default
		level: levels.Info,

		TraceOutput: levels.Trace.Output(), // os.Stderr,
		DebugOutput: levels.Debug.Output(), // os.Stderr,
		InfoOutput:  levels.Info.Output(),  // os.Stdout,
		WarnOutput:  levels.Warn.Output(),  // os.Stderr,
		ErrorOutput: levels.Error.Output(), // os.Stderr,
		PanicOutput: levels.Panic.Output(), // os.Stderr,
		FatalOutput: levels.Fatal.Output(), // os.Stderr,

		TraceColor: levels.Trace.Color(), // Green,
		DebugColor: levels.Debug.Color(), // Blue,
		WarnColor:  levels.Warn.Color(),  // Yellow,
		ErrorColor: levels.Error.Color(), // Red,
		PanicColor: levels.Panic.Color(), // RedB,
		FatalColor: levels.Fatal.Color(), // RedB,

	}
}

func (l *Logger) Named(name string) *Logger {
	var (
		colorTrace = *l.TraceColor
		colorDebug = *l.DebugColor
		colorWarn  = *l.WarnColor
		colorError = *l.ErrorColor
		colorPanic = *l.PanicColor
		colorFatal = *l.FatalColor
	)
	return &Logger{
		level:       l.level,
		prefix:      log.Prefixed(l.prefix, name),
		TraceOutput: l.TraceOutput,
		DebugOutput: l.DebugOutput,
		InfoOutput:  l.InfoOutput,
		WarnOutput:  l.WarnOutput,
		ErrorOutput: l.ErrorOutput,
		PanicOutput: l.PanicOutput,
		FatalOutput: l.FatalOutput,
		TraceColor:  &colorTrace,
		DebugColor:  &colorDebug,
		WarnColor:   &colorWarn,
		ErrorColor:  &colorError,
		PanicColor:  &colorPanic,
		FatalColor:  &colorFatal,
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

func (l Logger) log(level levels.Type, args ...interface{}) {
	if level < l.level { // Trace(0) < Info(2) => no logging
		return
	}
	var template string
	template, args = l.unmapped(args...)                   // fields in last position
	template, args = l.prefixedf(template, args...)        // prefix in first position
	template, args = l.levelledf(level, template, args...) // level in first position by pushing prefix to second
	msg := fmt.Sprintf(ln(template), args...)
	switch level {
	case levels.Trace:
		fmt.Fprint(l.TraceOutput, l.TraceColor.Sprint(msg))
	case levels.Debug:
		fmt.Fprint(l.DebugOutput, l.DebugColor.Sprint(msg))
	case levels.Warn:
		fmt.Fprint(l.WarnOutput, l.WarnColor.Sprint(msg))
	case levels.Error:
		fmt.Fprint(l.ErrorOutput, l.ErrorColor.Sprint(msg))
	case levels.Panic:
		fmt.Fprint(l.PanicOutput, l.PanicColor.Sprint(msg))
		panic(msg)
	case levels.Fatal:
		fmt.Fprint(l.FatalOutput, l.FatalColor.Sprint(msg))
		os.Exit(1)
	default: // levels.Info
		fmt.Fprint(l.InfoOutput, msg)
	}
}

func (l Logger) logf(level levels.Type, template string, args ...interface{}) {
	if level < l.level { // Trace(0) < Info(2) => no logging
		return
	}
	template, args = l.unmappedf(template, args...)        // fields values in last position
	template, args = l.prefixedf(template, args...)        // prefix in first position
	template, args = l.levelledf(level, template, args...) // level in first position by pushing prefix to second
	template = ln(template)
	msg := fmt.Sprintf(ln(template), args...)
	switch level {
	case levels.Trace:
		fmt.Fprint(l.TraceOutput, l.TraceColor.Sprint(msg))
	case levels.Debug:
		fmt.Fprint(l.DebugOutput, l.DebugColor.Sprint(msg))
	case levels.Warn:
		fmt.Fprint(l.WarnOutput, l.WarnColor.Sprint(msg))
	case levels.Error:
		fmt.Fprint(l.ErrorOutput, l.ErrorColor.Sprint(msg))
	case levels.Panic:
		fmt.Fprint(l.PanicOutput, l.PanicColor.Sprint(msg))
		panic(msg)
	case levels.Fatal:
		fmt.Fprint(l.FatalOutput, l.FatalColor.Sprint(msg))
		os.Exit(1)
	default: // levels.Info
		fmt.Fprint(l.InfoOutput, msg)
	}
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

func (l Logger) prefixedf(template string, args ...interface{}) (string, []interface{}) {
	if l.prefix != "" {
		return "%s: " + template, append([]interface{}{l.prefix}, args...)
	}
	return template, args
}

func (l Logger) levelledf(level levels.Type, template string, args ...interface{}) (string, []interface{}) {
	if level != levels.Info {
		return "%s " + template, append([]interface{}{level.Bracket()}, args...)
	}
	return template, args
}

func (l Logger) unmap(fields log.Map) string {
	var ret string
	for key, val := range fields {
		ret += fmt.Sprintf("%s=%s ", key, val)
	}
	return strings.TrimSuffix(ret, " ")
}

func ln(s string) string {
	if n := len(s); n == 0 || s[n-1] != '\n' {
		return s + "\n"
	}
	return s
}
