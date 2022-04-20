# Log Package

[log-go](https://github.com/Masterminds/log-go/blob/main/README.md) "spin-off" logger with its API closer to the Go standard library package.

[![Go Reference](https://pkg.go.dev/badge/github.com/roninzo/log.svg)](https://pkg.go.dev/github.com/roninzo/log)
[![license](https://img.shields.io/badge/license-MIT-green "The MIT License (MIT)")](LICENSE)
[![build](https://img.shields.io/badge/build-passing-green "Go build status")](log.go)
=======
<!-- [![coverage](https://img.shields.io/badge/coverage-65%25-orange?logo=codecov "Unit tests coverage")](example_test.go)  -->

# Usage

## Example Go code
```go
import (
	"fmt"

	"github.com/roninzo/log"
	"github.com/roninzo/log/impl/hclog"
	"github.com/roninzo/log/impl/logrus"
	"github.com/roninzo/log/impl/std"
	"github.com/roninzo/log/impl/zap"
)

type Foo struct {
	Logger log.Logger
}

func (f *Foo) DoSomething() {
	f.Logger.Info("Hello Logging")
	f.Logger.Info("Logging supports key-value pairs", log.Map{"foo": "bar"})
}

func main() {

	// Using the default Logger
	log.Info("Hello")
	log.Error("World")

	// Create a logrus logger with default configuration that uses the log
	// interface. Note, logrus can be setup with default settings or setup with
	// custom settings using a second constructor.
	lgrs := logrus.NewStandard()

	// Set logrus as the global logger
	log.Current = lgrs

	// Logrus is now used globally for logging
	log.Warn("Warning through logrus")

	f1 := Foo{
		Logger: lgrs,
	}

	// Logging in DoSomething will use the set logger which is logrus
	f1.DoSomething()

	f2 := Foo{
		// The log package uses the global logger from the standard library log
		// package. A custom standard library logger can be used with the
		// github.com/roninzo/log/impl/std package.
		Logger: log.NewStandard(),
	}

	// Logging in DoSomething will the logger from the standard library
	f2.DoSomething()

	// Need to detect the logger being used? You can check for the type.
	switch log.Current.(type) {
	case *log.Std:
		fmt.Println("The default logger")
	case *std.Logger:
		fmt.Println("The default logger")
	case *logrus.Logger:
		fmt.Printf("Logrus is used for logging")
	case *zap.Logger:
		fmt.Printf("Zap is used for logging")
	case *hclog.Logger:
		fmt.Printf("HashiCorp is used for logging")
	default:
		fmt.Printf("Something else that implements the interface")
	}
}

```

## Output
```log
2022/04/20 13:06:06 [INFO]  Hello
2022/04/20 13:06:06 [ERROR] World
time="2022-04-20T13:06:06+02:00" level=warning msg="Warning through logrus"
time="2022-04-20T13:06:06+02:00" level=info msg="Hello Logging"
time="2022-04-20T13:06:06+02:00" level=info msg="Logging supports key-value pairs" foo=bar
2022/04/20 13:06:06 [INFO]  Hello Logging
2022/04/20 13:06:06 [INFO]  Logging supports key-value pairs [foo=bar]
Logrus is used for logging
```

# Interface

## Exported logger interface
```go
type Logger interface {
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
```

## Implicit logger interface
Methods cannot be declared as part of the exported interface, but are available for all implementations.
```go
type Logger interface {
	Named(string) Logger                   // Create a new sub-Logger with a name descending from the current name. This is used to create a subsystem specific Logger.
	WithLevel(levels.Type) Logger          // Chainable level setter.
	WithLevelFromDebug(bool) Logger        // Chainable level setter from debug boolean value.
	Options(...func(Logger) Logger) Logger // Custom chainable setter functions.
}
```
