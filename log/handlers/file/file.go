package file

import (
	"bytes"
	"fmt"
	slog "log"
	"sort"
	"sync"

	"github.com/keiwi/utils/log"
	"github.com/kyokomi/emoji"
)

var (
	// FancyFormatter is the default formatter for fancy messages.
	FileFormatter = log.MustStringFormatter(`[{{formatTime .Timestamp "2006-01-02 15:04:05"}}] [{{ .ShortFile }}] {{ .LevelIcon }} {{ .LevelTitle }}{{"\t"}}> {{ .Message }}{{ .ParsedFields }}`)
)

func NewFile(config *Config) *File {
	w := writer{
		Mutex:      new(sync.Mutex),
		fileFormat: config.Filename,
		folder:     config.Folder,
		maxSize:    config.MaxSize,
		maxLines:   config.MaxLines,
	}

	err := w.Init()
	if err != nil {
		slog.Fatalf("Error when initializing file logger: %v", err.Error())
	}

	return &File{writer: w, Formatter: FileFormatter}
}

type Config struct {
	Filename string
	Folder   string
	MaxSize  int64
	MaxLines int64
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

	c.writer.Write([]byte(fmt.Sprintln(msg)))
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
