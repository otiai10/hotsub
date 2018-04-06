package core

import (
	"io"
)

// Report ...
type Report struct {
	Log     Logger
	Metrics struct {
		Writer io.Writer
	}
}

// Lifetime ...
func (job *Job) Lifetime(label, format string, v ...interface{}) {
	if job.Report.Log == nil {
		return
	}
	job.Report.Log.Lifetime(label, format, v...)
}

// Stdio logs stdout/stderr.
func (job *Job) Stdio(streamtype int, label string, text string) {
	if job.Report.Log == nil {
		return
	}
	job.Report.Log.Stdio(streamtype, label, text)

}
