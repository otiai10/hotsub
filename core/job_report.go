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

// Lifetimef ...
func (job *Job) Lifetimef(format string, v ...interface{}) {
	if job.Report.Log == nil {
		return
	}
	job.Report.Log.Lifetimef(format, v...)
}

// Stdf ...
func (job *Job) Stdf(streamtype int, format string, v ...interface{}) {
	if job.Report.Log == nil {
		return
	}
	job.Report.Log.Stdf(streamtype, format, v...)
}
