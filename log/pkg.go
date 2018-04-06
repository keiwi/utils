package log

import (
	"sync"

	"github.com/tryy3/joy/stdlib/strings"
)

// singletons ftw?
var Log = &Logger{
	Reporters: []Reporter{stdLog{}},
	Level:     INFO,
}

// AddReporter adds a new reporter to the logger
func AddReporter(r Reporter) {
	Log.Reporters = append(Log.Reporters, r)
}

// NewLogger creates a new logger
func NewLogger(level Level, reporters []Reporter) *Logger {
	return &Logger{
		Mutex:     new(sync.Mutex),
		Reporters: reporters,
		Level:     level,
	}
}

// SetLevel sets the log level. This is not thread-safe.
func SetLevel(l Level) {
	Log.Level = l
}

// SetLevelFromString sets the log level from a string, panicing when invalid. This is not thread-safe.
func SetLevelFromString(s string) {
	Log.Level = GetLevelFromString(s)
}

// GetLevelFromString returns the log level from a string, panicing when invalid
func GetLevelFromString(s string) Level {
	s = strings.ToLower(s)
	for _, level := range Levels {
		if s == strings.ToLower(level.Name) {
			return level.Level
		}
	}
	panic("Invalid level name (" + s + ")")
	return FATAL
}

// WithFields returns a new entry with `fields` set.
func WithFields(fields Fielder) *Entry {
	return NewEntry(Log).WithFields(fields)
}

// WithField returns a new entry with the `key` and `value` set.
func WithField(key string, value interface{}) *Entry {
	return NewEntry(Log).WithField(key, value)
}

// WithError returns a new entry with the "error" set to `err`.
func WithError(err error) *Entry {
	return NewEntry(Log).WithError(err)
}

// Debug level message.
func Debug(msg string) {
	NewEntry(Log).Debug(msg)
}

// Info level message.
func Info(msg string) {
	NewEntry(Log).Info(msg)
}

// Warn level message.
func Warn(msg string) {
	NewEntry(Log).Warn(msg)
}

// Error level message.
func Error(msg string) {
	NewEntry(Log).Error(msg)
}

// Fatal level message, followed by an exit.
func Fatal(msg string) {
	NewEntry(Log).Fatal(msg)
}

// Debugf level formatted message.
func Debugf(msg string, v ...interface{}) {
	NewEntry(Log).Debugf(msg, v...)
}

// Infof level formatted message.
func Infof(msg string, v ...interface{}) {
	NewEntry(Log).Infof(msg, v...)
}

// Warnf level formatted message.
func Warnf(msg string, v ...interface{}) {
	NewEntry(Log).Warnf(msg, v...)
}

// Errorf level formatted message.
func Errorf(msg string, v ...interface{}) {
	NewEntry(Log).Errorf(msg, v...)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...interface{}) {
	NewEntry(Log).Fatalf(msg, v...)
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func Trace(msg string) *Entry {
	return NewEntry(Log).Trace(msg)
}
