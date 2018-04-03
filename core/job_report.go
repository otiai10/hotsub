package core

import (
	"io"
	"strings"
)

// Report ...
type Report struct {
	Log     Logger
	Metrics struct {
		Writer io.Writer
	}
}

// Logf ...
func (job *Job) Logf(format string, v ...interface{}) {
	if job.Report.Log == nil {
		return
	}
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	job.Report.Log.Printf(format, v...)
}
