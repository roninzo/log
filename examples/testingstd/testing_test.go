package testingstd

import (
	"bytes"
	stdlog "log"
	"testing"

	"github.com/roninzo/log"
	"github.com/roninzo/log/impl/std"
	"github.com/roninzo/log/levels"
	"github.com/stretchr/testify/assert"
)

// This example is designed to help you understand how you can write tests with
// the interface and implementations.
// For further examples, each of the implementations in the impl directory has
// tests associated with you that check for expected output. You can use those
// as additional examples.
// This example uses the logger included in the standard library. The same idea
// and capability is included for the other loggers in the impl directory.

func TestLogger(t *testing.T) {

	// Create a logger that uses a buffer to capture the logging output
	var buf bytes.Buffer

	// The impl/std package enables you to pass in an instance of a logger from
	// the standard library. The github.com/roninzo/log.NewStandard()
	// constructor uses the standard logger in the standard libraries log pkg.
	logger := std.New(stdlog.New(&buf, "", stdlog.Lshortfile))

	// Set the appropriate log level for use in the tests. Info is the default.
	logger.Level(levels.Debug)

	// Try out the Fibonacci generator that uses the passed in logger. Ignoring
	// the result as we aren't testing the result.
	lgrfib(logger, 1)

	// See if the buffer the logger is writing to has the right response
	assert.Contains(t, buf.String(), `Number is 1`)

	// Empty the buffer before the next test.
	buf.Reset()

	// The package level logger can be tested as well.
	// First, save the current package level logger, setup a defer, and replace
	// it.
	prev := log.Current
	defer func() {
		log.Current = prev
	}()
	log.Current = logger

	// Try out the Fibonacci generator that uses package level logger
	fib(1)

	// See if the buffer the logger is writing to has the right response
	// In this case we check the log level that was written to the log
	assert.Contains(t, buf.String(), `[DEBUG] Number is 1`)

	// Empty the buffer before the next test.
	buf.Reset()

}

// A basic Fibonacci generator that logs the number passed in using the logger
// configured as the Current one for package log package level functions.
func fib(num uint) uint {
	log.Debugf("Number is %d", num)

	if num <= 1 {
		return num
	}

	return fib(num-1) + fib(num-2)
}

// A basic Fibonacci generator that logs the number passed in using the passed
// in logger that conforms to the log.Logger interface.
func lgrfib(logger log.Logger, num uint) uint {
	logger.Debugf("Number is %d", num)

	if num <= 1 {
		return num
	}

	return lgrfib(logger, num-1) + lgrfib(logger, num-2)
}
