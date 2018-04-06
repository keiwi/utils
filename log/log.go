package log

import (
	"sort"
	"sync"
	"time"

	stdlog "log"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// stackTracer interface.
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Fielder is an interface for providing fields to custom types.
type Fielder interface {
	Fields() Fields
}

// Fields represents a map of entry level data used for structured logging.
type Fields map[string]interface{}

// Fields implements Fielder.
func (f Fields) Fields() Fields {
	return f
}

// Get field value by name.
func (f Fields) Get(name string) interface{} {
	return f[name]
}

// Names returns field names sorted.
func (f Fields) Names() (v []string) {
	for k := range f {
		v = append(v, k)
	}

	sort.Strings(v)
	return
}

type Reporter interface {
	Write(e *Entry, calldepth int) error
}

type Logger struct {
	*sync.Mutex
	Level     Level
	Reporters []Reporter
}

//

// WithFields returns a new entry with `fields` set.
func (l *Logger) WithFields(fields Fielder) *Entry {
	return NewEntry(l).WithFields(fields.Fields())
}

// WithField returns a new entry with the `key` and `value` set.
//
// Note that the `key` should not have spaces in it - use camel
// case or underscores
func (l *Logger) WithField(key string, value interface{}) *Entry {
	return NewEntry(l).WithField(key, value)
}

// WithError returns a new entry with the "error" set to `err`.
func (l *Logger) WithError(err error) *Entry {
	return NewEntry(l).WithError(err)
}

// Debug level message.
func (l *Logger) Debug(msg string) {
	NewEntry(l).Debug(msg)
}

// Info level message.
func (l *Logger) Info(msg string) {
	NewEntry(l).Info(msg)
}

// Warn level message.
func (l *Logger) Warn(msg string) {
	NewEntry(l).Warn(msg)
}

// Error level message.
func (l *Logger) Error(msg string) {
	NewEntry(l).Error(msg)
}

// Fatal level message, followed by an exit.
func (l *Logger) Fatal(msg string) {
	NewEntry(l).Fatal(msg)
}

// Debugf level formatted message.
func (l *Logger) Debugf(msg string, v ...interface{}) {
	NewEntry(l).Debugf(msg, v...)
}

// Infof level formatted message.
func (l *Logger) Infof(msg string, v ...interface{}) {
	NewEntry(l).Infof(msg, v...)
}

// Warnf level formatted message.
func (l *Logger) Warnf(msg string, v ...interface{}) {
	NewEntry(l).Warnf(msg, v...)
}

// Errorf level formatted message.
func (l *Logger) Errorf(msg string, v ...interface{}) {
	NewEntry(l).Errorf(msg, v...)
}

// Fatalf level formatted message, followed by an exit.
func (l *Logger) Fatalf(msg string, v ...interface{}) {
	NewEntry(l).Fatalf(msg, v...)
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func (l *Logger) Trace(msg string) *Entry {
	return NewEntry(l).Trace(msg)
}

func (l *Logger) Write(level Level, e *Entry, msg string, calldepth int) *Logger {
	if e.Level > l.Level {
		return l
	}

	l.Lock()

	e.Timestamp = time.Now()
	finished := e.finalize(level, msg)

	var result error
	for _, r := range l.Reporters {
		if err := r.Write(finished, calldepth+1); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if result != nil {
		stdlog.Printf("error logging: %s", result)
	}
	l.Unlock()

	return l
}
