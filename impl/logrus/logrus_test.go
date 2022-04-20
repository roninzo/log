package logrus

import (
	"bytes"
	"testing"

	"github.com/roninzo/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogrus(t *testing.T) {

	// Test the logger meets the interface
	var _ log.Logger = new(Logger)

	var logger = logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	buf := &bytes.Buffer{}
	logger.SetOutput(buf)
	lgr := New(logger)

	lgr.Trace("test trace")
	assert.Contains(t, buf.String(), `level=trace msg="test trace"`)
	buf.Reset()

	lgr.Tracef("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=trace msg="Hello World"`)
	buf.Reset()

	lgr.Trace("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=trace msg="foo bar" baz=qux`)
	buf.Reset()

	lgr.Debug("test debug")
	assert.Contains(t, buf.String(), `level=debug msg="test debug"`)
	buf.Reset()

	lgr.Debugf("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=debug msg="Hello World"`)
	buf.Reset()

	lgr.Debug("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=debug msg="foo bar" baz=qux`)
	buf.Reset()

	lgr.Info("test info")
	assert.Contains(t, buf.String(), `level=info msg="test info"`)
	buf.Reset()

	lgr.Infof("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=info msg="Hello World"`)
	buf.Reset()

	lgr.Info("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=info msg="foo bar" baz=qux`)
	buf.Reset()

	lgr.Warn("test warn")
	assert.Contains(t, buf.String(), `level=warning msg="test warn"`)
	buf.Reset()

	lgr.Warnf("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=warning msg="Hello World"`)
	buf.Reset()

	lgr.Warn("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=warning msg="foo bar" baz=qux`)
	buf.Reset()

	lgr.Error("test error")
	assert.Contains(t, buf.String(), `level=error msg="test error"`)
	buf.Reset()

	lgr.Errorf("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=error msg="Hello World"`)
	buf.Reset()

	lgr.Error("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=error msg="foo bar" baz=qux`)
	buf.Reset()

	// lgr.Fatal("test fatal")
	// lgr.Fatalf(template string, args ...interface{})
	// lgr.Fatal(msg string, fields Map)

	assert.Panics(t, func() { lgr.Panic("test panic") })
	assert.Contains(t, buf.String(), `level=panic msg="test panic"`)
	buf.Reset()

	assert.Panics(t, func() { lgr.Panicf("Hello %s", "World") })
	assert.Contains(t, buf.String(), `level=panic msg="Hello World"`)
	buf.Reset()

	assert.Panics(t, func() { lgr.Panic("foo bar", log.Map{"baz": "qux"}) })
	assert.Contains(t, buf.String(), `level=panic msg="foo bar" baz=qux`)
	buf.Reset()
}

func TestStandardLogrus(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	buf := &bytes.Buffer{}
	logrus.SetOutput(buf)
	lgr := NewStandard()

	lgr.Trace("test trace")
	assert.NotContains(t, buf.String(), `level=trace msg="test trace"`)
	buf.Reset()

	lgr.Tracef("Hello %s", "World")
	assert.NotContains(t, buf.String(), `level=trace msg="Hello World"`)
	buf.Reset()

	lgr.Trace("foo bar", log.Map{"baz": "qux"})
	assert.NotContains(t, buf.String(), `level=trace msg="foo bar" baz=qux`)
	buf.Reset()

	lgr.Debug("test debug")
	assert.Contains(t, buf.String(), `level=debug msg="test debug"`)
	buf.Reset()

	lgr.Debugf("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=debug msg="Hello World"`)
	buf.Reset()

	lgr.Debug("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=debug msg="foo bar" baz=qux`)
	buf.Reset()

	lgr.Info("test info")
	assert.Contains(t, buf.String(), `level=info msg="test info"`)
	buf.Reset()

	lgr.Infof("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=info msg="Hello World"`)
	buf.Reset()

	lgr.Info("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=info msg="foo bar" baz=qux`)
	buf.Reset()

	lgr.Warn("test warn")
	assert.Contains(t, buf.String(), `level=warning msg="test warn"`)
	buf.Reset()

	lgr.Warnf("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=warning msg="Hello World"`)
	buf.Reset()

	lgr.Warn("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=warning msg="foo bar" baz=qux`)
	buf.Reset()

	lgr.Error("test error")
	assert.Contains(t, buf.String(), `level=error msg="test error"`)
	buf.Reset()

	lgr.Errorf("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=error msg="Hello World"`)
	buf.Reset()

	lgr.Error("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=error msg="foo bar" baz=qux`)
	buf.Reset()

	// lgr.Fatal("test fatal")
	// lgr.Fatalf(template string, args ...interface{})
	// lgr.Fatal(msg string, fields Map)

	assert.Panics(t, func() { lgr.Panic("test panic") })
	assert.Contains(t, buf.String(), `level=panic msg="test panic"`)
	buf.Reset()

	assert.Panics(t, func() { lgr.Panicf("Hello %s", "World") })
	assert.Contains(t, buf.String(), `level=panic msg="Hello World"`)
	buf.Reset()

	assert.Panics(t, func() { lgr.Panic("foo bar", log.Map{"baz": "qux"}) })
	assert.Contains(t, buf.String(), `level=panic msg="foo bar" baz=qux`)
	buf.Reset()
}

func TestLogrusInterface(t *testing.T) {
	lgr := NewStandard()
	testfunc(lgr)
}

func testfunc(l log.Logger) {
	l.Debug("test")
}
