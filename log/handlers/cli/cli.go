package cli

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/keiwi/utils/log"
	"github.com/kyokomi/emoji"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
)

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

func NewCli() *Cli {
	return &Cli{colorable.NewColorableStdout(), CliFormatter}
}

var (
	// FancyFormatter is the default formatter for fancy messages.
	CliFormatter = log.MustStringFormatter(`[{{formatColorString "magentafg" (formatTime .Timestamp "2006-01-02 15:04:05")}}] [{{formatColorString "magentafg" .ShortFile}}] {{formatColor .LevelColor .LevelIcon .LevelTitle "\t>"}} {{ .Message }}{{ .ParsedFields }}`)
)

type Cli struct {
	Writer    io.Writer
	Formatter log.Formatter
}

func (c *Cli) Write(e *log.Entry, calldepth int) error {
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

	_, err = fmt.Fprintln(colorable.NewColorableStdout(), msg)
	return err
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
			field = aurora.Colorize(field, aurora.BlackFg).Bold().String()
		} else {
			field = aurora.Colorize(field, aurora.BlackFg).Bold().String()
		}
		if len(fields) > 2 {
			field = "\n\t" + field
		} else {
			field = aurora.Colorize("| ", aurora.BlackFg).Bold().String() + field
		}
		fmt.Fprintf(&b, " %s", field)
	}
	if len(fields) > 2 {
		fmt.Fprintln(&b, "")
	}

	return b.String()
}
