package std

import (
	"bytes"
	stdlog "log"
	"testing"

	"github.com/roninzo/log"
	"github.com/roninzo/log/levels"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {

	// Test the logger meets the interface
	var _ log.Logger = new(Logger)

	buf := &bytes.Buffer{}
	logger := stdlog.New(buf, "MyPrefix: ", stdlog.Lshortfile)
	lgr := New(logger)

	// Make sure levels are working
	lgr.Debug("test debug")
	assert.NotContains(t, buf.String(), `[DEBUG] test debug`)
	buf.Reset()

	// Test all levels
	lgr.Level(levels.Trace)

	lgr.Trace("test trace")
	assert.Contains(t, buf.String(), `[TRACE] test trace`)
	buf.Reset()

	lgr.Tracef("Hello %s", "World")
	assert.Contains(t, buf.String(), `[TRACE] Hello World`)
	buf.Reset()

	lgr.Trace("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `[TRACE] foo bar [baz=qux]`)
	buf.Reset()

	lgr.Debug("test debug")
	assert.Contains(t, buf.String(), `[DEBUG] test debug`)
	buf.Reset()

	lgr.Debugf("Hello %s", "World")
	assert.Contains(t, buf.String(), `[DEBUG] Hello World`)
	buf.Reset()

	lgr.Debug("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `[DEBUG] foo bar [baz=qux]`)
	buf.Reset()

	lgr.Info("test info")
	assert.Contains(t, buf.String(), `[INFO]  test info`)
	buf.Reset()

	lgr.Infof("Hello %s", "World")
	assert.Contains(t, buf.String(), `[INFO]  Hello World`)
	buf.Reset()

	lgr.Info("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `[INFO]  foo bar [baz=qux]`)
	buf.Reset()

	lgr.Warn("test warn")
	assert.Contains(t, buf.String(), `[WARN]  test warn`)
	buf.Reset()

	lgr.Warnf("Hello %s", "World")
	assert.Contains(t, buf.String(), `[WARN]  Hello World`)
	buf.Reset()

	lgr.Warn("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `[WARN]  foo bar [baz=qux]`)
	buf.Reset()

	lgr.Error("test error")
	assert.Contains(t, buf.String(), `[ERROR] test error`)
	buf.Reset()

	lgr.Errorf("Hello %s", "World")
	assert.Contains(t, buf.String(), `[ERROR] Hello World`)
	buf.Reset()

	lgr.Error("foo bar", log.Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `[ERROR] foo bar [baz=qux]`)
	buf.Reset()

	// lgr.Fatal("test fatal")
	// lgr.Fatalf(template string, args ...interface{})
	// lgr.Fatal(msg string, fields Map)

	assert.PanicsWithValue(t, "[PANIC] test panic\n", func() { lgr.Panic("test panic") })
	assert.Contains(t, buf.String(), `test panic`)
	buf.Reset()

	assert.PanicsWithValue(t, "[PANIC] Hello World\n", func() { lgr.Panicf("Hello %s", "World") })
	assert.Contains(t, buf.String(), `Hello World`)
	buf.Reset()

	assert.PanicsWithValue(t, "[PANIC] foo bar [baz=qux]\n", func() { lgr.Panic("foo bar", log.Map{"baz": "qux"}) })
	assert.Contains(t, buf.String(), `foo bar [baz=qux]`)
	buf.Reset()
}
