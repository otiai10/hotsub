package logs

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sync"

	"github.com/fatih/color"
	"github.com/otiai10/awsub/core"
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

type (
	// ColorfulLoggerFactory ...
	ColorfulLoggerFactory struct {
		w io.Writer
	}
	// ColorfulLogger ...
	ColorfulLogger struct {
		w io.Writer
		c *color.Color
		p string
	}
)

// Logger ...
func (clf *ColorfulLoggerFactory) Logger(job *core.Job) core.Logger {
	if clf.w == nil {
		clf.w = os.Stdout
	}
	return &ColorfulLogger{
		w: clf.w,
		c: colors[job.Identity.Index%len(colors)],
		p: fmt.Sprintf("[%s %d]", job.Identity.Prefix, job.Identity.Index),
	}
}

// Printf ...
func (logger *ColorfulLogger) Printf(format string, v ...interface{}) {

	// Force newline at the tail
	format = newline.ReplaceAllString(format, "\n")

	// TODO: Refactor
	// To avoid to mix up prefix and log content, lock the print process for each print.
	linelock.Lock()
	defer linelock.Unlock()

	logger.c.Print(logger.p + " ")
	fmt.Printf(format, v...)
}
