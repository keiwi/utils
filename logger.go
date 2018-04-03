package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	stdlog "log"

	"github.com/apex/log"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/afero"
)

var Log *log.Logger

func init() {
	Log = NewLogger(log.DebugLevel, &LoggerConfig{})
}

var NameToLevel = map[string]log.Level{
	"debug": log.DebugLevel,
	"info":  log.InfoLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
	"fatal": log.FatalLevel,
}

var levelNames = [...]string{
	log.DebugLevel: "DEBUG",
	log.InfoLevel:  "INFO ",
	log.WarnLevel:  "WARN ",
	log.ErrorLevel: "ERROR",
	log.FatalLevel: "FATAL",
}

var levelColors = [...]aurora.Color{
	log.DebugLevel: aurora.MagentaFg,
	log.InfoLevel:  aurora.GreenFg,
	log.WarnLevel:  aurora.BrownFg,
	log.ErrorLevel: aurora.RedFg,
	log.FatalLevel: aurora.CyanBg,
}

func OpenFile(fs afero.Fs, path string, flag int, perm os.FileMode) (afero.File, error) {
	if ok, err := afero.Exists(fs, path); !ok || err != nil {
		dir := filepath.Dir(path)
		err = fs.MkdirAll(dir, perm)
		if err != nil {
			return nil, err
		}
	}

	return fs.OpenFile(path, flag, perm)
}

type LoggerConfig struct {
	Dirname string
	Logname string
}

type Logger struct {
	Config *LoggerConfig
	fs     afero.Fs
	lock   sync.Mutex
	queue  []*log.Entry
}

func (l *Logger) HandleLog(e *log.Entry) error {
	var callstack string
	if Log.Level == log.DebugLevel {
		pc := make([]uintptr, 10)
		n := runtime.Callers(1, pc)
		if n != 0 {
			pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
			frames := runtime.CallersFrames(pc)

			// Loop to get frames.
			// A fixed number of pcs can expand to an indefinite number of Frames.
			for {
				frame, more := frames.Next()
				if !strings.Contains(frame.File, "github.com/keiwi/utils/logger.go") &&
					!strings.Contains(frame.File, "github.com/apex/log") {
					file := frame.File[strings.Index(frame.File, "src")+4:]
					if callstack == "" {
						callstack = fmt.Sprintf("%s:%d", file, frame.Line)
					} else {
						callstack = fmt.Sprintf("%s:%d -> %s", file, frame.Line, callstack)
					}
				}
				if !more {
					break
				}
			}
		}
		e.Fields["callstack"] = callstack
	}

	l.println(e)
	go l.addLog(e)
	return nil
}

func (l *Logger) addLog(e *log.Entry) {
	l.lock.Lock()
	l.queue = append(l.queue, e)
	l.lock.Unlock()

	time.Sleep(time.Second * 1)
	l.handle()
}

func (l *Logger) handle() {
	l.lock.Lock()
	defer l.lock.Unlock()
	if len(l.queue) <= 0 {
		return
	}

	if l.Config.Dirname == "" || l.Config.Dirname == "" {
		return
	}

	if err := l.writeToLogFile(); err != nil {
		stdlog.Printf("error logging: %s", err)
	}

	l.queue = []*log.Entry{}
}

func (l *Logger) println(e *log.Entry) {
	fmt.Println(l.formatEntry(e, true))
}

func (l *Logger) writeToLogFile() error {
	file := l.Config.Logname
	date := time.Now().Format("2006-01-02")
	file = strings.Replace(file, "%date%", date, -1)
	path := filepath.Join(l.Config.Dirname, file)

	f, err := OpenFile(l.fs, path, os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	for _, e := range l.queue {
		data := l.formatEntry(e, false)
		_, err = f.WriteString(data + "\n")
	}
	return err
}

func (l *Logger) formatEntry(e *log.Entry, color bool) string {
	level := levelNames[e.Level]
	date := e.Timestamp.Format("2006-01-02 15:04:05")

	format := "[%s] %s -\t%s"
	if color {
		c := levelColors[e.Level]
		return fmt.Sprintf(
			format,
			aurora.Colorize(date, aurora.MagentaFg).Bold(),
			aurora.Colorize(level, c),
			aurora.Colorize(l.parseEntry(e), c),
		)
	}

	return fmt.Sprintf(
		format,
		date,
		level,
		l.parseEntry(e),
	)
}

// field used for sorting.
type field struct {
	Name  string
	Value interface{}
}

// by sorts fields by name.
type byName []field

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (l *Logger) parseEntry(e *log.Entry) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s", e.Message)

	var fields []field
	var callstack *field

	for k, v := range e.Fields {
		if k == "callstack" {
			callstack = &field{k, v}
			continue
		}
		fields = append(fields, field{k, v})
	}

	sort.Sort(byName(fields))
	if callstack != nil {
		fields = append(fields, *callstack)
	}

	for _, v := range fields {
		field := fmt.Sprintf("%s=%v", v.Name, v.Value)
		if v.Name == "callstack" {
			field = aurora.Colorize(field, aurora.BlackFg).Bold().String()
		} else {
			field = aurora.Colorize(field, aurora.MagentaFg).String()
		}
		field = aurora.Colorize("| ", aurora.MagentaFg).String() + field
		fmt.Fprintf(&b, " %s", field)
	}

	return b.String()
}

func NewLogger(level log.Level, config *LoggerConfig) *log.Logger {
	return &log.Logger{
		Level: level,
		Handler: &Logger{
			Config: config,
			fs:     afero.NewOsFs(),
		},
	}
}
