package log

import (
	"bytes"
	"fmt"
	"log"
	"sort"
)

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

type stdLog struct{}

func (stdLog) Write(e *Entry, calldepth int) error {
	level := Levels[e.Level]

	var fields []field

	for k, v := range e.Fields {
		fields = append(fields, field{k, v})
	}

	sort.Sort(byName(fields))

	var b bytes.Buffer
	fmt.Fprintf(&b, "%5s %-25s", level, e.Message)

	for _, f := range fields {
		fmt.Fprintf(&b, " %s=%v", f.Name, f.Value)
	}

	log.Println(b.String())

	return nil
}
