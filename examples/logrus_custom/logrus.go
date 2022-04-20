package main

import (
	"github.com/roninzo/log"
	lgrs "github.com/roninzo/log/impl/logrus"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	log.Current = lgrs.New(logger)

	log.Debug("Debugging!")

	// A basic message
	log.Info("Hello,", "World")

	// Use Go formatting on a warning
	log.Warnf("Foo %s", "bar")

	// An error with context
	log.Error("foo, bar", log.Map{"baz": "qux"})
}
