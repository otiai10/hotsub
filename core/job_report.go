package core

import (
	"io"

	"github.com/otiai10/daap"
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
func (job *Job) Stdio(streamtype daap.HijackedStreamType, label Lifecycle, text string) {
	if job.Report.Log == nil {
		return
	}
	job.Report.Log.Stdio(int(streamtype), string(label), text)

}
