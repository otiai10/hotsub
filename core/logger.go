package core

// Loggerer FIXME: Should it be moved to anywhere else?
type Loggerer interface {
	Logger(*Job) Logger
}

// Logger ...
type Logger interface {
	Printf(string, ...interface{})
}
