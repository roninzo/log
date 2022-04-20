// Package io provides a means of turning a log.Logger into an io.Writer for
// a chosen level. An example of this would be:
//     import(
//         "io"
//         "github.com/roninzo/log"
//         logio "github.com/roninzo/log/io"
//     )
//
//     func main() {
//         w := logio.NewCurrentWriter(levels.Info)
//         io.WriteString(w, "foo")
//     }
package io

import (
	"github.com/roninzo/log"
	"github.com/roninzo/log/levels"
)

// CurrentWriter uses the current package level logger for io writing
type CurrentWriter struct {
	Level levels.Type
}

// NewCurrentWriter creates a new CurrentWriter. The levels that can be passed
// to it are:
// - levels.Trace:
// - levels.Debug:
// - levels.Info:
// - levels.Warn:
// - levels.Error:
// - levels.Panic:
// - levels.Fatal:
func NewCurrentWriter(level levels.Type) *CurrentWriter {
	return &CurrentWriter{
		Level: level,
	}
}

// Write is the write method from the io.Writer interface in the standard lib
func (l CurrentWriter) Write(p []byte) (n int, err error) {
	switch l.Level {
	case levels.Trace:
		log.Trace(string(p))
	case levels.Debug:
		log.Debug(string(p))
	case levels.Info:
		log.Info(string(p))
	case levels.Warn:
		log.Warn(string(p))
	case levels.Error:
		log.Error(string(p))
	case levels.Panic:
		log.Panic(string(p))
	case levels.Fatal:
		log.Fatal(string(p))
	default:
		log.Panicf("Invalid logger level selected: %d", l.Level)
	}

	return len(p), nil
}

// Writer uses the configured logger for io writing
type Writer struct {
	Logger log.Logger
	Level  levels.Type
}

// NewWriter creates a new Writer. It accepts a logger and a level that
// will be written on the io.Writer interface. The levels you can pass in are:
// - levels.Trace:
// - levels.Debug:
// - levels.Info:
// - levels.Warn:
// - levels.Error:
// - levels.Panic:
// - levels.Fatal:
func NewWriter(lgr log.Logger, level levels.Type) *Writer {
	return &Writer{
		Logger: lgr,
		Level:  level,
	}
}

// Write is the write method from the io.Writer interface in the standard lib
func (l Writer) Write(p []byte) (n int, err error) {
	switch l.Level {
	case levels.Trace:
		l.Logger.Trace(string(p))
	case levels.Debug:
		l.Logger.Debug(string(p))
	case levels.Info:
		l.Logger.Info(string(p))
	case levels.Warn:
		l.Logger.Warn(string(p))
	case levels.Error:
		l.Logger.Error(string(p))
	case levels.Panic:
		l.Logger.Panic(string(p))
	case levels.Fatal:
		l.Logger.Fatal(string(p))
	default:
		l.Logger.Panicf("Invalid logger level selected: %d", l.Level)
	}

	return len(p), nil
}
