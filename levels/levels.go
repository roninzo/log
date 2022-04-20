package levels

import (
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Type int

const (
	Trace Type = iota
	Debug
	Info
	Warn
	Error
	Panic
	Fatal
	Silent
)

// FromString returns a Level type for the named log level, or "Undef"
// if the level string is invalid. This facilitates setting the log level
// via config or environment variable by name in a predictable way.
// Example:
// 	"debug" => level.Debug
// The Level string argument `s` is NOT "case-sensitive", i.e.:
// 	Accept both "INFO" and "info".
func FromString(s string) Type {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "trace":
		return Trace
	case "debug":
		return Debug
	case "info":
		return Info
	case "warn", "warning":
		return Warn
	case "error":
		return Error
	case "panic":
		return Panic
	case "fatal":
		return Fatal
	case "silent", "off":
		return Silent
	default:
		return Silent
	}
}

// FromDebug returns the debug level type if debug argument is true,
// or the info level type if debug argument if false or missing.
// This facilitates setting the log level via config or environment
// variable by value in a predictable way.
// Example:
// 	true      => level.Debug
// 	false     => level.Info
// 	<nothing> => level.Info
func FromDebug(debug ...bool) Type {
	if len(debug) > 0 && debug[0] {
		return Debug
	}
	return Info
}

func (l Type) String() string {
	switch l {
	case Trace:
		return "trace"
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warn:
		return "warning"
	case Error:
		return "error"
	case Panic:
		return "panic"
	case Fatal:
		return "fatal"
	case Silent:
		return "silent"
	default:
		return "undef"
	}
}

func (l Type) Bracket() string {
	switch l {
	case Trace:
		return "[TRACE]"
	case Debug:
		return "[DEBUG]"
	case Info:
		return "[INFO] "
	case Warn:
		return "[WARN] "
	case Error:
		return "[ERROR]"
	case Panic:
		return "[PANIC]"
	case Fatal:
		return "[FATAL]"
	case Silent:
		return ""
	default:
		return "[UNDEF]"
	}
}

func (l Type) Output() io.Writer {
	switch l {
	case Info:
		return os.Stdout
	default:
		return os.Stderr
	}
}

func (l Type) Color() *color.Color {
	switch l {
	case Trace:
		return color.New(color.FgHiGreen)
	case Debug:
		return color.New(color.FgHiWhite)
	case Info:
		return color.New(color.FgHiBlue)
	case Warn:
		return color.New(color.FgHiYellow)
	case Error:
		return color.New(color.FgHiRed)
	case Panic:
		return color.New(color.FgHiRed, color.Bold)
	case Fatal:
		return color.New(color.FgHiRed, color.Bold)
	default:
		return color.New(color.FgWhite)
	}
}
