package zap

import (
	"fmt"
	"strings"
	"testing"

	"github.com/roninzo/log"
	"github.com/roninzo/log/levels"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestLogger(t *testing.T) {

	// Test the logger meets the interface
	var _ log.Logger = new(Logger)

	ts := newTestLogSpy(t)
	defer ts.AssertPassed()
	atom := zap.NewAtomicLevelAt(zap.DebugLevel - 1) // ~ levels.Trace
	// logger := zaptest.NewLogger(ts, zapcore.NewCore(atom))
	logger := zaptest.NewLogger(ts, zaptest.Level(atom.Level()))
	defer func() {
		err := logger.Sync()
		if err != nil {
			t.Errorf("Error syncing logger: %s", err)
		}
	}()

	// // TODO: Fix the following code and decomission fatal
	// // ending.
	// lgr := New(logger)
	lgr := New(logger, atom)

	assert.Equal(t, levels.Trace, lgr.getLevel())            // Trace != Fatal
	assert.Equal(t, atom.Level(), lgr.intLevel(lgr.Level())) // -2 != 5

	lgr.Trace("test trace")
	lgr.Tracef("Hello %s", "World")
	lgr.Trace("foo bar", log.Map{"baz": "qux"})
	lgr.Debug("test debug")
	lgr.Debugf("Hello %s", "World")
	lgr.Debug("foo bar", log.Map{"baz": "qux"})
	lgr.Info("test info")
	lgr.Infof("Hello %s", "World")
	lgr.Info("foo bar", log.Map{"baz": "qux"})
	lgr.Warn("test warn")
	lgr.Warnf("Hello %s", "World")
	lgr.Warn("foo bar", log.Map{"baz": "qux"})
	lgr.Error("test error")
	lgr.Errorf("Hello %s", "World")
	lgr.Error("foo bar", log.Map{"baz": "qux"})

	assert.Panics(t, func() { lgr.Panic("test panic") })
	assert.Panics(t, func() { lgr.Panicf("Hello %s", "World") })
	assert.Panics(t, func() { lgr.Panic("foo bar", log.Map{"baz": "qux"}) })

	ts.AssertMessages(
		"LEVEL(-2)	test trace",
		"LEVEL(-2)	Hello World",
		"LEVEL(-2)	foo bar	{\"baz\": \"qux\"}",
		"DEBUG	test debug",
		"DEBUG	Hello World",
		"DEBUG	foo bar	{\"baz\": \"qux\"}",
		"INFO	test info",
		"INFO	Hello World",
		"INFO	foo bar	{\"baz\": \"qux\"}",
		"WARN	test warn",
		"WARN	Hello World",
		"WARN	foo bar	{\"baz\": \"qux\"}",
		"ERROR	test error",
		"ERROR	Hello World",
		"ERROR	foo bar	{\"baz\": \"qux\"}",
		"PANIC	test panic",
		"PANIC	Hello World",
		"PANIC	foo bar	{\"baz\": \"qux\"}",
	)

	// lgr.Fatal("test fatal")
	// lgr.Fatalf(template string, args ...interface{})
	// lgr.Fatal(msg string, fields Map)
}

func TestZapInterface(t *testing.T) {
	logger := zaptest.NewLogger(t, zaptest.Level(zap.DebugLevel))
	defer func() {
		err := logger.Sync()
		if err != nil {
			t.Errorf("Error syncing logger: %s", err)
		}
	}()
	lgr := New(logger)
	testfunc(lgr)
}

func testfunc(l log.Logger) {
	l.Debug("test")
}

// This last section is taken from zap itself and licensed under the MIT license
type testLogSpy struct {
	testing.TB

	failed   bool
	Messages []string
}

func newTestLogSpy(t testing.TB) *testLogSpy {
	return &testLogSpy{TB: t}
}

func (t *testLogSpy) Fail() {
	t.failed = true
}

func (t *testLogSpy) Failed() bool {
	return t.failed
}

func (t *testLogSpy) FailNow() {
	t.Fail()
	t.TB.FailNow()
}

func (t *testLogSpy) Logf(format string, args ...interface{}) {
	// Log messages are in the format,
	//
	//   2017-10-27T13:03:01.000-0700	DEBUG	your message here	{data here}
	//
	// We strip the first part of these messages because we can't really test
	// for the timestamp from these tests.
	m := fmt.Sprintf(format, args...)
	m = m[strings.IndexByte(m, '\t')+1:]
	t.Messages = append(t.Messages, m)
	t.TB.Log(m)
}

func (t *testLogSpy) AssertMessages(msgs ...string) {
	assert.Equal(t.TB, msgs, t.Messages, "logged messages did not match")
}

func (t *testLogSpy) AssertPassed() {
	t.assertFailed(false, "expected test to pass")
}

func (t *testLogSpy) AssertFailed() {
	t.assertFailed(true, "expected test to fail")
}

func (t *testLogSpy) assertFailed(v bool, msg string) {
	assert.Equal(t.TB, v, t.failed, msg)
}
