package main

import (
	"github.com/roninzo/log"
	"github.com/roninzo/log/impl/logrus"
)

func main() {
	log.Current = logrus.NewStandard()

	// A basic message
	log.Info("Hello,", "World")

	// Use Go formatting on a warning
	log.Warnf("Foo %s", "bar")

	// An error with context
	log.Error("foo, bar", log.Map{"baz": "qux"})
}
