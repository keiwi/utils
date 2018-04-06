package log

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/logrusorgru/aurora"
)

var (
	// DefaultFormatter is the default formatter used and is only the message.
	DefaultFormatter = MustStringFormatter("{{ .Message }}")

	// FancyFormatter is the default formatter for fancy messages.
	FancyFormatter = MustStringFormatter(`[{{formatColorString "magentafg" (formatTime .Timestamp "2006-01-02 15:04:05")}}] {{formatColor .LevelColor .LevelTitle "\t>"}} {{ .Message }}`)
)

var (
	pid     = os.Getpid()
	program = filepath.Base(os.Args[0])
)

type Formatter interface {
	Format(*Entry, int) map[string]interface{}
	Finalize(map[string]interface{}) (string, error)
}

type formatter struct {
	format *template.Template
}

func (f *formatter) Format(e *Entry, calldepth int) map[string]interface{} {
	calldepth = calldepth + 1
	pc, file, line, ok := runtime.Caller(calldepth)

	// Short and long file processing
	shortf := file
	if !ok {
		file = "???"
		shortf = file
		line = 0
	} else {
		shortf = filepath.Base(file)
	}
	shortfile := fmt.Sprintf("%s:%d", filepath.Base(shortf), line)
	longfile := fmt.Sprintf("%s:%d", file, line)

	// Package and func name processing
	longpkg := "???"
	shortpkg := "???"
	longfunc := "???"
	shortfunc := "???"
	if ok {
		if f := runtime.FuncForPC(pc); f != nil {
			longpkg = formatFuncName("longpkg", f.Name())
			shortpkg = formatFuncName("shortpkg", f.Name())
			longfunc = formatFuncName("longfunc", f.Name())
			shortfunc = formatFuncName("shortfunc", f.Name())
		}
	}

	format := map[string]interface{}{
		"PID":        pid,
		"Timestamp":  e.Timestamp,
		"Level":      e.Level,
		"LevelTitle": Levels[e.Level].Name,
		"LevelColor": Levels[e.Level].Color,
		"Program":    program,
		"Message":    e.Message,
		"LongFile":   longfile,
		"ShortFile":  shortfile,
		"LongPkg":    longpkg,
		"ShortPkg":   shortpkg,
		"LongFunc":   longfunc,
		"ShortFunc":  shortfunc,
		"Calldepth":  calldepth,
	}
	return format
}

func (f *formatter) Finalize(data map[string]interface{}) (string, error) {
	var b bytes.Buffer
	err := f.format.ExecuteTemplate(&b, "formatter", data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func NewStringFormatter(format string) (Formatter, error) {
	tmpl := template.New("formatter")

	funcs := template.FuncMap{
		"formatCallpath":    formatCallpath,
		"formatColor":       formatColor,
		"formatColorString": formatColorString,
		"formatTime":        formatTime,
	}

	tmpl = tmpl.Funcs(funcs)

	tmpl, err := tmpl.Parse(format)
	if err != nil {
		return nil, err
	}
	return &formatter{tmpl}, nil
}

// MustStringFormatter is equivalent to NewStringFormatter with a call to panic
// on error.
func MustStringFormatter(format string) Formatter {
	f, err := NewStringFormatter(format)
	if err != nil {
		panic("Failed to initialized formatter: " + err.Error())
	}
	return f
}

// formatFuncName tries to extract certain part of the runtime formatted
// function name to some pre-defined variation.
//
// This function is known to not work properly if the package path or name
// contains a dot.
func formatFuncName(v string, f string) string {
	i := strings.LastIndex(f, "/")
	j := strings.Index(f[i+1:], ".")
	if j < 1 {
		return "???"
	}
	pkg, fun := f[:i+j+1], f[i+j+2:]
	switch v {
	case "longpkg":
		return pkg
	case "shortpkg":
		return path.Base(pkg)
	case "longfunc":
		return fun
	case "shortfunc":
		i = strings.LastIndex(fun, ".")
		return fun[i+1:]
	}
	panic("unexpected func formatter")
}

func formatCallpath(calldepth int, depth int) string {
	v := ""
	callers := make([]uintptr, 64)
	n := runtime.Callers(calldepth+2, callers)
	oldPc := callers[n-1]

	start := n - 3
	if depth > 0 && start >= depth {
		start = depth - 1
		v += "~."
	}
	recursiveCall := false
	for i := start; i >= 0; i-- {
		pc := callers[i]
		if oldPc == pc {
			recursiveCall = true
			continue
		}
		oldPc = pc
		if recursiveCall {
			recursiveCall = false
			v += ".."
		}
		if i < start {
			v += "."
		}
		if f := runtime.FuncForPC(pc); f != nil {
			v += formatFuncName("shortfunc", f.Name())
		}
	}
	return v
}

func formatColor(color aurora.Color, v ...interface{}) string {
	return aurora.Colorize(fmt.Sprint(v...), color).String()
}

func formatColorString(color string, v ...interface{}) string {
	var c aurora.Color
	switch strings.ToLower(color) {
	case "blackfg":
		c = aurora.BlackFg
	case "redfg":
		c = aurora.RedFg
	case "greenfg":
		c = aurora.GreenFg
	case "brownfg":
		c = aurora.BrownFg
	case "bluefg":
		c = aurora.BlueFg
	case "magentafg":
		c = aurora.MagentaFg
	case "cyanfg":
		c = aurora.CyanFg
	case "grayfg":
		c = aurora.GrayFg
	case "blackbg":
		c = aurora.BlackBg
	case "redbg":
		c = aurora.RedBg
	case "greenbg":
		c = aurora.GreenBg
	case "brownbg":
		c = aurora.BrownBg
	case "bluebg":
		c = aurora.BlueBg
	case "magentabg":
		c = aurora.MagentaBg
	case "cyanbg":
		c = aurora.CyanBg
	case "graybg":
		c = aurora.GrayBg
	case "bold":
		c = aurora.BoldFm
	case "inverse":
		c = aurora.InverseFm
	}
	return formatColor(c, v...)
}

func formatTime(t time.Time, format string) string {
	return t.Format(format)
}
