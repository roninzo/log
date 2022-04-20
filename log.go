// ⚡️ The log package provides a common interface that can be used in applications and libraries along with reference implementation wrappers for logrus, zap, the Go standard library package, and for a CLI.
// 🤖 Forked Github Repository: github.com/Masterminds/log-go
// 📌 API Documentation: https://github.com/Masterminds/log-go/blob/main/README.md

package log

import "github.com/roninzo/log/levels"

// Logger is an interface for Logging.
type Logger interface {
	//                                     NOTE: Methods cannot be declared as part of the exported interface, but are available for all implementations.
	// Named(string) Logger                   // Create a new sub-Logger with a name descending from the current name. This is used to create a subsystem specific Logger.
	// WithLevel(levels.Type) Logger          // Chainable level setter.
	// WithLevelFromDebug(bool) Logger        // Chainable level setter from debug boolean value.
	// Options(...func(Logger) Logger) Logger // Custom chainable setter functions.
	Prefix(...string) string          // Prefix returns current logger name. With a prefix argument, the current logger's name is set to it.
	Level(...levels.Type) levels.Type // Level returns current logging level. With a level argument, the current logger's level is set to it.
	Trace(...interface{})             // Trace logs a message at the Trace level.
	Debug(...interface{})             // Debug logs a message at the Debug level.
	Info(...interface{})              // Info logs a message at the Info level.
	Warn(...interface{})              // Warn logs a message at the Warn level.
	Error(...interface{})             // Error logs a message at the Error level.
	Panic(...interface{})             // Panic logs a message at the Panic level and panics.
	Fatal(...interface{})             // Fatal logs a message at the Fatal level and exists the application.
	Tracef(string, ...interface{})    // Tracef formats a message according to a format specifier and logs the message at the Trace level.
	Debugf(string, ...interface{})    // Debugf formats a message according to a format specifier and logs the message at the Debug level.
	Infof(string, ...interface{})     // Infof formats a message according to a format specifier and logs the message at the Info level.
	Warnf(string, ...interface{})     // Warnf formats a message according to a format specifier and logs the message at the Warning level.
	Errorf(string, ...interface{})    // Errorf formats a message according to a format specifier and logs the message at the Error level.
	Panicf(string, ...interface{})    // Panicf formats a message according to a format specifier and logs the message at the Panic level and then panics.
	Fatalf(string, ...interface{})    // Fatalf formats a message according to a format specifier and logs the message at the Fatal level and exits the application.
}

// Current contains the logger used for the package level logging functions.
var Current Logger

func init() {
	Current = NewStandard()
}

func Prefix(prefix ...string) string              { return Current.Prefix(prefix...) }
func Level(level ...levels.Type) levels.Type      { return Current.Level(level...) }
func Trace(msg ...interface{})                    { Current.Trace(msg...) }
func Debug(msg ...interface{})                    { Current.Debug(msg...) }
func Info(msg ...interface{})                     { Current.Info(msg...) }
func Warn(msg ...interface{})                     { Current.Warn(msg...) }
func Error(msg ...interface{})                    { Current.Error(msg...) }
func Panic(msg ...interface{})                    { Current.Panic(msg...) }
func Fatal(msg ...interface{})                    { Current.Fatal(msg...) }
func Tracef(template string, args ...interface{}) { Current.Tracef(template, args...) }
func Debugf(template string, args ...interface{}) { Current.Debugf(template, args...) }
func Infof(template string, args ...interface{})  { Current.Infof(template, args...) }
func Warnf(template string, args ...interface{})  { Current.Warnf(template, args...) }
func Errorf(template string, args ...interface{}) { Current.Errorf(template, args...) }
func Panicf(template string, args ...interface{}) { Current.Panicf(template, args...) }
func Fatalf(template string, args ...interface{}) { Current.Fatalf(template, args...) }
