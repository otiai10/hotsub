package core

// Loggerer can generate a Logger struct corresponding to the specified Job.
type Loggerer interface {
	Logger(*Job) Logger
}

// Logger is an interface of log writer
type Logger interface {
	Lifetimef(string, ...interface{})
	Stdf(int, string, ...interface{})
	Close() error
}
