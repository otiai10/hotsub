package main

// Job represents an instance of what handles a task, and it's logs and results.
type Job struct {
	Task  Task
	Error error
}
