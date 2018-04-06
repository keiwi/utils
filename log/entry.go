package log

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Entry represents a single log entry.
type Entry struct {
	Logger    *Logger
	Level     Level
	Message   string
	Formatted map[string]interface{}
	Fields    Fields
	Timestamp time.Time
	start     time.Time
	fields    []Fields
	calldepth int
}

// NewEntry returns a new entry for `log`.
func NewEntry(log *Logger) *Entry {
	return &Entry{
		Logger:    log,
		calldepth: 1,
	}
}

// WithFields returns a new entry with `fields` set.
func (e *Entry) WithFields(fields Fielder) *Entry {
	f := []Fields{}
	f = append(f, e.fields...)
	f = append(f, fields.Fields())
	return &Entry{
		Logger:    e.Logger,
		fields:    f,
		calldepth: e.calldepth - 1,
	}
}

// WithField returns a new entry with the `key` and `value` set.
func (e *Entry) WithField(key string, value interface{}) *Entry {
	return e.WithFields(Fields{key: value})
}

// WithError returns a new entry with the "error" set to `err`.
//
// The given error may implement .Fielder, if it does the method
// will add all its `.Fields()` into the returned entry.
func (e *Entry) WithError(err error) *Entry {
	//e.calldepth = e.calldepth + 1
	ctx := e.WithField("error", err.Error())

	if s, ok := err.(stackTracer); ok {
		frame := s.StackTrace()[0]

		name := fmt.Sprintf("%n", frame)
		file := fmt.Sprintf("%+s", frame)
		line := fmt.Sprintf("%d", frame)

		parts := strings.Split(file, "\n\t")
		if len(parts) > 1 {
			file = parts[1]
		}

		ctx = ctx.WithField("source", fmt.Sprintf("%s: %s:%s", name, file, line))
	}

	if f, ok := err.(Fielder); ok {
		ctx = ctx.WithFields(f.Fields())
	}

	return ctx
}

// Fatal level message, followed by an exit.
func (e *Entry) Fatal(msg string) {
	e.Logger.Write(FATAL, e, msg, e.calldepth+1)
	os.Exit(1)
}

// Error level message.
func (e *Entry) Error(msg string) {
	e.Logger.Write(ERROR, e, msg, e.calldepth+1)
}

// Warn level message.
func (e *Entry) Warn(msg string) {
	e.Logger.Write(WARN, e, msg, e.calldepth+1)
}

// Info level message.
func (e *Entry) Info(msg string) {
	e.Logger.Write(INFO, e, msg, e.calldepth+1)
}

// Debug level message.
func (e *Entry) Debug(msg string) {
	e.Logger.Write(DEBUG, e, msg, e.calldepth+1)
}

// Fatalf level formatted message, followed by an exit.
func (e *Entry) Fatalf(msg string, v ...interface{}) {
	e.calldepth = e.calldepth + 1
	e.Fatal(fmt.Sprintf(msg, v...))
}

// Errorf level formatted message.
func (e *Entry) Errorf(msg string, v ...interface{}) {
	e.calldepth = e.calldepth + 1
	e.Error(fmt.Sprintf(msg, v...))
}

// Warnf level formatted message.
func (e *Entry) Warnf(msg string, v ...interface{}) {
	e.calldepth = e.calldepth + 1
	e.Warn(fmt.Sprintf(msg, v...))
}

// Infof level formatted message.
func (e *Entry) Infof(msg string, v ...interface{}) {
	e.calldepth = e.calldepth + 1
	e.Info(fmt.Sprintf(msg, v...))
}

// Debugf level formatted message.
func (e *Entry) Debugf(msg string, v ...interface{}) {
	e.calldepth = e.calldepth + 1
	e.Debug(fmt.Sprintf(msg, v...))
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion write, useful with defer.
func (e *Entry) Trace(msg string) *Entry {
	e.calldepth = e.calldepth + 1
	e.Info(msg)
	v := e.WithFields(e.Fields)
	v.Message = msg
	v.start = time.Now()
	return v
}

// Stop should be used with Trace, to fire off the completion message. When
// an `err` is passed the "error" field is set, and the write level is error.
func (e *Entry) Stop(err *error) {
	e.calldepth = e.calldepth + 1
	if err == nil || *err == nil {
		e.WithField("duration", time.Since(e.start)).Info(e.Message)
	} else {
		e.WithField("duration", time.Since(e.start)).WithError(*err).Error(e.Message)
	}
}

// mergedFields returns the fields list collapsed into a single map.
func (e *Entry) mergedFields() Fields {
	f := Fields{}

	for _, fields := range e.fields {
		for k, v := range fields {
			f[k] = v
		}
	}

	return f
}

// finalize returns a copy of the Entry with Fields merged.
func (e *Entry) finalize(level Level, msg string) *Entry {
	entry := &Entry{
		Logger:    e.Logger,
		Fields:    e.mergedFields(),
		Level:     level,
		Message:   msg,
		Timestamp: time.Now(),
	}
	return entry
}
