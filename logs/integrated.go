package logs

import (
	"fmt"

	"github.com/otiai10/hotsub/core"
)

type (
	// IntegratedLoggerFactory ...
	IntegratedLoggerFactory struct {
		File    *FileLoggerFactory
		Verbose *ColorfulLoggerFactory
	}
	// IntegratedLogger ...
	IntegratedLogger struct {
		File    core.Logger
		Verbose core.Logger
	}
)

// Logger satisfies core.LoggerFactory
func (factory *IntegratedLoggerFactory) Logger(job *core.Job) (core.Logger, error) {
	logger := new(IntegratedLogger)
	if factory.File == nil && factory.Verbose == nil {
		return nil, fmt.Errorf("either of file logger or verebose logger must be specified")
	}
	if factory.File != nil {
		f, err := factory.File.Logger(job)
		if err != nil {
			return nil, err
		}
		logger.File = f
	}
	if factory.Verbose != nil {
		v, err := factory.Verbose.Logger(job)
		if err != nil {
			return nil, err
		}
		logger.Verbose = v
	}
	return logger, nil
}

// Lifetime ...
func (logger *IntegratedLogger) Lifetime(lifecycle string, format string, v ...interface{}) {
	if logger.File != nil {
		logger.File.Lifetime(lifecycle, format, v...)
	}
	if logger.Verbose != nil {
		logger.Verbose.Lifetime(lifecycle, format, v...)
	}
}

// Stdio ...
func (logger *IntegratedLogger) Stdio(streamtype int, lifecycle string, text string) {
	if logger.File != nil {
		logger.File.Stdio(streamtype, lifecycle, text)
	}
	if logger.Verbose != nil {
		logger.Verbose.Stdio(streamtype, lifecycle, text)
	}

}

// Close ...
func (logger *IntegratedLogger) Close() error {
	if logger.File != nil {
		if err := logger.File.Close(); err != nil {
			return err
		}
	}
	if logger.Verbose != nil {
		if err := logger.Verbose.Close(); err != nil {
			return err
		}
	}
	return nil
}
