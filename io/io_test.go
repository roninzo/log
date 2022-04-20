package io

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/roninzo/log"
	"github.com/roninzo/log/levels"
	"github.com/stretchr/testify/assert"
)

func TestCurrentIO(t *testing.T) {
	buf := &bytes.Buffer{}
	lgr := newTestLogger(buf)
	def := log.Current
	log.Current = lgr
	defer func() {
		log.Current = def
	}()

	o := NewCurrentWriter(levels.Trace)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=trace msg="testing"`)
	buf.Reset()

	o = NewCurrentWriter(levels.Debug)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=debug msg="testing"`)
	buf.Reset()

	o = NewCurrentWriter(levels.Info)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=info msg="testing"`)
	buf.Reset()

	o = NewCurrentWriter(levels.Warn)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=warning msg="testing"`)
	buf.Reset()

	o = NewCurrentWriter(levels.Error)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=error msg="testing"`)
	buf.Reset()

	o = NewCurrentWriter(levels.Panic)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=panic msg="testing"`)
	buf.Reset()

	o = NewCurrentWriter(levels.Fatal)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=fatal msg="testing"`)
	buf.Reset()

	// Testing a non-existant level
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
		}
	}()
	log.Current = def // Test logger does not panic. Using one that does
	o = NewCurrentWriter(5000)
	_, _ = io.WriteString(o, "not happening")
	assert.Contains(t, buf.String(), `level=000 msg="testing"`)
	buf.Reset()
	t.Error("Not happening Logger failed to panic")
}

func TestIO(t *testing.T) {
	buf := &bytes.Buffer{}
	lgr := newTestLogger(buf)

	o := NewWriter(lgr, levels.Trace)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=trace msg="testing"`)
	buf.Reset()

	o = NewWriter(lgr, levels.Debug)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=debug msg="testing"`)
	buf.Reset()

	o = NewWriter(lgr, levels.Info)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=info msg="testing"`)
	buf.Reset()

	o = NewWriter(lgr, levels.Warn)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=warning msg="testing"`)
	buf.Reset()

	o = NewWriter(lgr, levels.Error)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=error msg="testing"`)
	buf.Reset()

	o = NewWriter(lgr, levels.Panic)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=panic msg="testing"`)
	buf.Reset()

	o = NewWriter(lgr, levels.Fatal)
	_, _ = io.WriteString(o, "testing")
	assert.Contains(t, buf.String(), `level=fatal msg="testing"`)
	buf.Reset()

	// Testing a non-existant level
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
		}
	}()
	o = NewWriter(log.Current, 5000) // Test logger does not panic. Using one that does
	_, _ = io.WriteString(o, "not happening")
	assert.Contains(t, buf.String(), `level=000 msg="testing"`)
	buf.Reset()
	t.Error("Not happening Logger failed to panic")
}

// logt is a test Logger.
type logt struct {
	logger *bytes.Buffer
	level  levels.Type
	prefix string
}

const (
	defaultTemplate = "msg=%q"
)

func newTestLogger(buf *bytes.Buffer) *logt {
	return &logt{
		logger: buf,
		level:  levels.Info,
	}
}

func (l *logt) Named(name string) *logt {
	return &logt{
		logger: &bytes.Buffer{},
		level:  l.level,
		prefix: log.Prefixed(l.prefix, name),
	}
}

func (l *logt) Options(funcs ...func(*logt) *logt) *logt {
	for _, f := range funcs {
		f(l)
	}
	return l
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

// Helpers.

func (l logt) log(level levels.Type, msg ...interface{}) {
	l.level = level
	template, msg := l.unmapped(msg...)           // fields in last position
	template, msg = l.prefixedf(template, msg...) // prefix in first position
	template, msg = l.levelledf(template, msg...) // level in first position by pushing prefix to second
	l.logger.WriteString(ln(fmt.Sprintf(template, msg...)))
}

func (l logt) logf(level levels.Type, template string, args ...interface{}) {
	l.level = level
	template, args = l.unmappedf(template, args...) // fields values in last position
	template, args = l.prefixedf(template, args...) // prefix in first position
	template, args = l.levelledf(template, args...) // level in first position by pushing prefix to second
	l.logger.WriteString(ln(fmt.Sprintf(template, args...)))
}

func (l logt) unmapped(msg ...interface{}) (string, []interface{}) {
	var fields log.Map
	msg, fields = log.ParseArgs(msg...)
	msg = []interface{}{fmt.Sprint(msg...)}
	template := defaultTemplate
	if len(fields) > 0 {
		return template + " %s", append(msg, l.unmap(fields))
	}
	return template, msg
}

func (l logt) unmappedf(template string, args ...interface{}) (string, []interface{}) {
	var fields log.Map
	args, fields = log.ParseArgs(args...)
	msg := []interface{}{fmt.Sprintf(template, args...)}
	template = defaultTemplate
	if len(fields) > 0 {
		return template + " %s", append(msg, l.unmap(fields))
	}
	return template, msg
}

func (l logt) prefixedf(template string, args ...interface{}) (string, []interface{}) {
	if l.prefix != "" {
		return "prefix=%s " + template, append([]interface{}{l.prefix}, args...)
	}
	return template, args
}

func (l logt) levelledf(template string, args ...interface{}) (string, []interface{}) {
	return "level=%s " + template, append([]interface{}{l.level.String()}, args...)
}

func (l logt) unmap(fields log.Map) string {
	var ret string
	for key, val := range fields {
		switch t := val.(type) {
		case []byte:
			ret += fmt.Sprintf("%s=%q ", key, string(t))
		case string:
			ret += fmt.Sprintf("%s=%q ", key, t)
		default: // numerical or supports Stringer interface
			ret += fmt.Sprintf("%s=%s ", key, t)
		}
	}
	return strings.TrimSuffix(ret, " ")
}

func ln(s string) string {
	// // Trimming an extra space at the end of a logging line is not worth the computing
	// // effort. That is why the following code is commented out.
	// if n := len(s); n == 0 || s[n-1] != '\n' {
	// 	return s + "\n"
	// }
	return s
}
