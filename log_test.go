package log

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/roninzo/log/levels"
	"github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	buf := &bytes.Buffer{}
	lgr := newTestLogger(buf)
	func1 := func(l *logt) *logt {
		return l.Named("roninzo")
	}
	func2 := func(l *logt) *logt {
		return l.WithLevelFromDebug(false)
	}
	func3 := func(l *logt) *logt {
		return l.WithLevel(levels.Trace)
	}
	Current = lgr.Options(func1, func2, func3)
	// => append "(.)roninzo" to logger prefix, set level to Info, and last, set level to Trace.

	Trace("test trace")
	assert.Contains(t, buf.String(), `level=trace msg="test trace"`)
	buf.Reset()

	Tracef("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=trace msg="Hello World"`)
	buf.Reset()

	Trace("foo bar", Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=trace msg="foo bar" baz=qux`)
	buf.Reset()

	Debug("test debug")
	assert.Contains(t, buf.String(), `level=debug msg="test debug"`)
	buf.Reset()

	Debugf("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=debug msg="Hello World"`)
	buf.Reset()

	Debug("foo bar", Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=debug msg="foo bar" baz=qux`)
	buf.Reset()

	Info("test info")
	assert.Contains(t, buf.String(), `level=info msg="test info"`)
	buf.Reset()

	Infof("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=info msg="Hello World"`)
	buf.Reset()

	Info("foo bar", Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=info msg="foo bar" baz=qux`)
	buf.Reset()

	Warn("test warn")
	assert.Contains(t, buf.String(), `level=warning msg="test warn"`)
	buf.Reset()

	Warnf("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=warning msg="Hello World"`)
	buf.Reset()

	Warn("foo bar", Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=warning msg="foo bar" baz=qux`)
	buf.Reset()

	Error("test error")
	assert.Contains(t, buf.String(), `level=error msg="test error"`)
	buf.Reset()

	Errorf("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=error msg="Hello World"`)
	buf.Reset()

	Error("foo bar", Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=error msg="foo bar" baz=qux`)
	buf.Reset()

	Panic("test panic")
	assert.Contains(t, buf.String(), `level=panic msg="test panic"`)
	buf.Reset()

	Panicf("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=panic msg="Hello World"`)
	buf.Reset()

	Panic("foo bar", Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=panic msg="foo bar" baz=qux`)
	buf.Reset()

	Fatal("test fatal")
	assert.Contains(t, buf.String(), `level=fatal msg="test fatal"`)
	buf.Reset()

	Fatalf("Hello %s", "World")
	assert.Contains(t, buf.String(), `level=fatal msg="Hello World"`)
	buf.Reset()

	Fatal("foo bar", Map{"baz": "qux"})
	assert.Contains(t, buf.String(), `level=fatal msg="foo bar" baz=qux`)
	buf.Reset()
}

// logt is a test Logger.
type logt struct {
	logger *bytes.Buffer
	level  levels.Type
	prefix string
}

func newTestLogger(lgr *bytes.Buffer) *logt {
	return &logt{
		logger: lgr,
	}
}

func (l *logt) Options(funcs ...func(*logt) *logt) *logt {
	for _, f := range funcs {
		f(l)
	}
	return l
}

func (l *logt) Named(name string) *logt {
	return &logt{
		logger: &bytes.Buffer{},
		prefix: Prefixed(l.prefix, name),
		level:  l.level,
	}
}

func (l *logt) WithLevel(level levels.Type) *logt {
	l.Level(level)
	return l
}

func (l *logt) WithLevelFromDebug(debug bool) *logt {
	switch debug {
	case true:
		l.Level(levels.Debug)
	default:
		l.Level(levels.Info)
	}
	return l
}

func (l *logt) Prefix(prefix ...string) string {
	if len(prefix) > 0 {
		l.prefix = prefix[0]
	}
	return l.prefix
}

func (l *logt) Level(level ...levels.Type) levels.Type {
	if len(level) > 0 {
		l.level = level[0]
	}
	return l.level
}

func (l logt) Trace(msg ...interface{}) { l.log(levels.Trace, msg...) }
func (l logt) Debug(msg ...interface{}) { l.log(levels.Debug, msg...) }
func (l logt) Info(msg ...interface{})  { l.log(levels.Info, msg...) }
func (l logt) Warn(msg ...interface{})  { l.log(levels.Warn, msg...) }
func (l logt) Error(msg ...interface{}) { l.log(levels.Error, msg...) }
func (l logt) Panic(msg ...interface{}) { l.log(levels.Panic, msg...) }
func (l logt) Fatal(msg ...interface{}) { l.log(levels.Fatal, msg...) }

func (l logt) Tracef(template string, args ...interface{}) { l.logf(levels.Trace, template, args...) }
func (l logt) Debugf(template string, args ...interface{}) { l.logf(levels.Debug, template, args...) }
func (l logt) Infof(template string, args ...interface{})  { l.logf(levels.Info, template, args...) }
func (l logt) Warnf(template string, args ...interface{})  { l.logf(levels.Warn, template, args...) }
func (l logt) Errorf(template string, args ...interface{}) { l.logf(levels.Error, template, args...) }
func (l logt) Panicf(template string, args ...interface{}) { l.logf(levels.Panic, template, args...) }
func (l logt) Fatalf(template string, args ...interface{}) { l.logf(levels.Fatal, template, args...) }

func (l logt) log(level levels.Type, msg ...interface{}) {
	template, args := l.unmapped(msg...)                   // fields in last position
	template, args = l.prefixedf(template, args...)        // prefix in first position
	template, args = l.levelledf(level, template, args...) // level in first position by pushing prefix to second
	l.logger.WriteString(ln(fmt.Sprintf(template, args...)))
}

func (l logt) logf(level levels.Type, template string, args ...interface{}) {
	template, args = l.unmappedf(template, args...)        // fields in last position
	template, args = l.prefixedf(template, args...)        // prefix in first position
	template, args = l.levelledf(level, template, args...) // level in first position by pushing prefix to second
	l.logger.WriteString(ln(fmt.Sprintf(template, args...)))
}

func (l logt) prefixedf(template string, args ...interface{}) (string, []interface{}) {
	if l.prefix != "" {
		return "prefix=%s " + template, append([]interface{}{l.prefix}, args...)
	}
	return template, args
}

func (l logt) levelledf(level levels.Type, template string, args ...interface{}) (string, []interface{}) {
	return "level=%s " + template, append([]interface{}{level.String()}, args...)
}

func (l logt) unmapped(args ...interface{}) (string, []interface{}) {
	var fields Map
	args, fields = ParseArgs(args...)
	args = []interface{}{fmt.Sprint(args...)}
	template := defaultTemplate
	if len(fields) > 0 {
		return template + " %s", append(args, l.unmap(fields))
	}
	return template, args
}

func (l logt) unmappedf(template string, args ...interface{}) (string, []interface{}) {
	var fields Map
	args, fields = ParseArgs(args...)
	args = []interface{}{fmt.Sprintf(template, args...)}
	template = defaultTemplate
	if len(fields) > 0 {
		return template + " %s", append(args, l.unmap(fields))
	}
	return template, args
}

func (l logt) unmap(fields Map) string {
	var ret string
	for key, val := range fields {
		ret += fmt.Sprintf("%s=%s ", key, val)
	}
	return strings.TrimSuffix(ret, " ")
}
