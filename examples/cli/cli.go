package main

import (
	"github.com/roninzo/log"
	"github.com/roninzo/log/impl/cli"
	"github.com/roninzo/log/levels"
)

func main() {
	logger := cli.NewStandard()
	logger.Level(levels.Trace) // logger.Level = log.Trace

	log.Current = logger

	// A basic message
	log.Info("Hello,", "World")

	// A trace message
	log.Trace("A low level trace message")

	// A trace message
	log.Debugf("Hello, %s", "World")

	// Use Go formatting on a warning
	log.Warnf("Foo %s", "bar")

	// An error with context
	log.Error("foo, bar", log.Map{"baz": "qux"})
}
