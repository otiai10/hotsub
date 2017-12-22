package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var colors = []*color.Color{
	color.New(color.FgHiCyan),
	color.New(color.FgHiGreen),
	color.New(color.FgHiBlue),
	color.New(color.FgHiMagenta),
	color.New(color.FgHiYellow),
	color.New(color.FgCyan),
	color.New(color.FgGreen),
	color.New(color.FgBlue),
	color.New(color.FgMagenta),
	color.New(color.FgYellow),
}

// Logger ...
type Logger struct {
	prefix string
	color  *color.Color
}

// NewLogger ...
func NewLogger(prefix string, index int) *Logger {
	return &Logger{
		prefix: prefix,
		color:  colors[index%len(colors)],
	}
}

// Printf ...
func (l *Logger) Printf(format string, v ...interface{}) {
	l.color.Print(l.prefix + " ")
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	fmt.Printf(format, v...)
}
