package log

import (
	"github.com/logrusorgru/aurora"
)

type Level int

const (
	FATAL Level = iota
	ERROR
	WARN
	INFO
	DEBUG
)

type LevelInfo struct {
	Level Level
	Color aurora.Color
	Name  string
}

var Levels = map[Level]LevelInfo{
	FATAL: {Level: FATAL, Color: aurora.RedFg, Name: "Fatal"},
	ERROR: {Level: ERROR, Color: aurora.RedFg, Name: "Error"},
	WARN:  {Level: WARN, Color: aurora.BrownFg, Name: "Warn"},
	INFO:  {Level: INFO, Color: aurora.CyanFg, Name: "Info"},
	DEBUG: {Level: DEBUG, Color: aurora.MagentaFg, Name: "Debug"},
}
