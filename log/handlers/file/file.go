package file

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/keiwi/utils/log"
	"github.com/kyokomi/emoji"
)

type writer struct {
	*sync.Mutex
	format   string
	folder   string
	messages []string
}

func (w *writer) log(message string) {
	w.Lock()
	w.messages = append(w.messages, message)
	w.Unlock()

	time.Sleep(time.Second * 1)
	w.write()
}

func (w *writer) write() {
	w.Lock()
	if len(w.messages) < 0 {
		w.Unlock()
		return
	}

	file := w.format
	date := time.Now().Format("2006-01-02")
	file = strings.Replace(file, "%date%", date, -1)
	path := filepath.Join(w.folder, file)

	f, err := w.open(path, os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		w.Unlock()
		panic(err.Error())
		return
	}

	for _, e := range w.messages {
		_, err = f.WriteString(e)
		if err != nil {
			f.Close()
			w.Unlock()
			panic(err.Error())
			return
		}
	}
	w.messages = []string{}

	f.Close()
	w.Unlock()
	return
}

func (w *writer) open(file string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(file, flag, perm)
}

var icons = map[log.Level]string{
	log.FATAL: ":shit:",
	log.ERROR: ":no_entry_sign:",
	log.WARN:  ":x:",
	log.INFO:  ":+1:",
	log.DEBUG: ":mag:",
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

func NewFile(folder string, format string) *File {
	w := writer{new(sync.Mutex), format, folder, []string{}}
	return &File{w, FileFormatter}
}

var (
	// FancyFormatter is the default formatter for fancy messages.
	FileFormatter = log.MustStringFormatter(`[{{formatTime .Timestamp "2006-01-02 15:04:05"}}] [{{ .ShortFile }}] {{ .LevelIcon }} {{ .LevelTitle }}{{"\t"}}> {{ .Message }}{{ .ParsedFields }}`)
)

type File struct {
	writer    writer
	Formatter log.Formatter
}

func (c *File) Write(e *log.Entry, calldepth int) error {
	fields := make([]field, len(e.Fields))
	i := 0
	for k, v := range e.Fields {
		fields[i] = field{k, v}
		i++
	}
	sort.Sort(byName(fields))

	data := c.Formatter.Format(e, calldepth+1)
	data["LevelIcon"] = emoji.Sprintf(icons[e.Level])
	data["Fields"] = fields
	data["ParsedFields"] = parseEntry(e)

	e.Formatted = data

	msg, err := c.Formatter.Finalize(e.Formatted)
	if err != nil {
		return err
	}

	go c.writer.log(fmt.Sprintln(msg))
	return nil
}

func parseEntry(e *log.Entry) string {
	var b bytes.Buffer

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
			field = field
		} else {
			field = field
		}
		if len(fields) > 2 {
			field = "\n\t" + field
		} else {
			field = "| " + field
		}
		fmt.Fprintf(&b, " %s", field)
	}
	if len(fields) > 2 {
		fmt.Fprintln(&b, "")
	}

	return b.String()
}
