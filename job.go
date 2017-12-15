package main

import (
	"fmt"
	"log"

	"github.com/otiai10/dkmachine/v0/dkmachine"
)

// Job represents an instance of what handles a task, and it's logs and results.
type Job struct {
	Task     *Task
	Error    error
	Instance *dkmachine.CreateOptions
	Logger   *log.Logger `json:"-"`
}

// Errorf is a shorthand for setting error and rerutn job.
func (job *Job) Errorf(format string, v ...interface{}) *Job {
	job.Error = fmt.Errorf(format, v...)
	return job
}

// Logf ...
func (job *Job) Logf(format string, v ...interface{}) {
	if job.Logger == nil {
		return
	}
	job.Logger.Printf(format, v...)
}
