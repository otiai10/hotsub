package main

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/fatih/color"
)

var colors = []*color.Color{
	color.New(color.FgHiGreen),
	color.New(color.FgHiCyan),
	color.New(color.FgHiMagenta),
	color.New(color.FgHiBlue),
	color.New(color.FgHiYellow),
	color.New(color.FgCyan),
	color.New(color.FgGreen),
	color.New(color.FgBlue),
	color.New(color.FgMagenta),
	color.New(color.FgYellow),
}

var (
	// TODO: Refactor
	linelock = new(sync.Mutex)

	// log suffix
	newline = regexp.MustCompile("\n*$")
)

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

	// if !strings.HasSuffix(format, "\n") {
	// 	format = format + "\n"
	// }
	format = newline.ReplaceAllString(format, "\n")

	// TODO: Refactor
	// To avoid to mix up prefix and log content, lock the print process for each print.
	linelock.Lock()
	defer linelock.Unlock()

	l.color.Print(l.prefix + " ")
	fmt.Printf(format, v...)
}
