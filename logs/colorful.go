package logs

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
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

// Lifetime output log of the job lifetime.
func (logger *ColorfulLogger) Lifetime(label, format string, v ...interface{}) {
	format = fmt.Sprintf("[%s]\t", strings.ToUpper(label)) + format
	logger.printf(format, v...)
}

// Stdio logs to appropriate writer according to streamtype [stdout, stderr]
func (logger *ColorfulLogger) Stdio(streamtype int, label, text string) {
	out := fmt.Sprintf("[%s]\t&%d> %s", strings.ToUpper(label), streamtype, text)
	logger.printf(out)
}

// Close ...
func (logger *ColorfulLogger) Close() error {
	if closer, ok := logger.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
