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
		writer io.Writer
		color  *color.Color
		prefix string
	}
)

// Logger ...
func (clf *ColorfulLoggerFactory) Logger(job *core.Job) core.Logger {
	if clf.w == nil {
		clf.w = os.Stdout
	}
	return &ColorfulLogger{
		writer: clf.w,
		color:  colors[job.Identity.Index%len(colors)],
		prefix: fmt.Sprintf("[%s %d]", job.Identity.Prefix, job.Identity.Index),
	}
}

// printf ...
func (logger *ColorfulLogger) printf(format string, v ...interface{}) {

	// Force newline at the tail
	format = newline.ReplaceAllString(format, "\n")

	// TODO: Refactor
	// To avoid to mix up prefix and log content, lock the print process for each print.
	linelock.Lock()
	defer linelock.Unlock()

	logger.color.Print(logger.prefix + "\t")
	fmt.Printf(format, v...)
}

// Lifetimef output log of the job lifetime.
func (logger *ColorfulLogger) Lifetimef(format string, v ...interface{}) {
	logger.printf(format, v...)
}

// Stdf output log to appropriate writer according to streamtype [stdout, stderr]
func (logger *ColorfulLogger) Stdf(streamtype int, format string, v ...interface{}) {
	logger.printf(format, v...)
}

// Close ...
func (logger *ColorfulLogger) Close() error {
	return nil
}
