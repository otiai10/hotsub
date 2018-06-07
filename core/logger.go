package core

// LoggerFactory can generate a Logger struct corresponding to the specified Job.
type LoggerFactory interface {
	Logger(*Job) (Logger, error)
}

// Logger is an interface of log writer
type Logger interface {
	Lifetime(string, string, ...interface{})
	Stdio(int, string, string)
	Close() error
}

func (component *Component) loggerForJob(job *Job) error {
	if component.JobLoggerFactory == nil {
		return nil
	}
	logger, err := component.JobLoggerFactory.Logger(job)
	if err != nil {
		return err
	}
	job.Report.Log = logger
	return nil
}
