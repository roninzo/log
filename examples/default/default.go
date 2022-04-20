package main

import "github.com/roninzo/log"

func main() {

	// A basic message
	log.Info("Hello,", "World")

	// Use Go formatting on a warning
	log.Warnf("Foo %s", "bar")

	// An error with context
	log.Error("foo, bar", log.Map{"baz": "qux"})
}
