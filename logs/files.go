package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/otiai10/awsub/core"
)

type (
	// FileLoggerFactory can construct FileLogger.
	FileLoggerFactory struct {
		Dir string
	}
	// FileLogger can dispatch and write logs to corresponding log files.
	FileLogger struct {
		rootdir  string
		stdout   *os.File
		stderr   *os.File
		lifetime *os.File
	}
)

// Logger constructs FileLogger, which satisfies core.Logger interface.
func (factory *FileLoggerFactory) Logger(job *core.Job) (core.Logger, error) {

	root, err := filepath.Abs(factory.Dir)
	if err != nil {
		return nil, err
	}
	jobdir := filepath.Join(root, job.Identity.Name)

	if err := os.MkdirAll(jobdir, os.ModePerm); err != nil {
		return nil, err
	}

	stdout, err := os.Create(filepath.Join(jobdir, "stdout.log"))
	if err != nil {
		return nil, err
	}

	stderr, err := os.Create(filepath.Join(jobdir, "stderr.log"))
	if err != nil {
		return nil, err
	}

	lifetime, err := os.Create(filepath.Join(jobdir, "lifetime.log"))
	if err != nil {
		return nil, err
	}

	log.Printf("[COMMAND]\tSee logs -> %s\n", jobdir)

	return &FileLogger{
		stdout:   stdout,
		stderr:   stderr,
		lifetime: lifetime,
	}, nil
}

// fprintf ...
func (logger *FileLogger) fprintf(w io.Writer, format string, v ...interface{}) error {
	// Force newline at the tail
	format = newline.ReplaceAllString(format, "\n")
	_, err := fmt.Fprintf(w, format, v...)
	return err
}

// Lifetime ...
func (logger *FileLogger) Lifetime(label, format string, v ...interface{}) {
	format = time.Now().Format(time.RFC3339) + "\t" + label + "\t" + format
	logger.fprintf(logger.lifetime, format, v...)
}

// Stdio ...
func (logger *FileLogger) Stdio(streamtype int, label string, text string) {
	switch streamtype {
	case 1:
		logger.fprintf(logger.stdout, "%s\t%s", label, text)
	case 2:
		logger.fprintf(logger.stderr, "%s\t%s", label, text)
	default:
		logger.fprintf(logger.lifetime, "%s\t?%d>\t%s", label, streamtype, text)
	}
}

// Close ...
func (logger *FileLogger) Close() error {
	if err := logger.closeIfExists(logger.stdout); err != nil {
		return err
	}
	if err := logger.closeIfExists(logger.stderr); err != nil {
		return err
	}
	if err := logger.closeIfExists(logger.lifetime); err != nil {
		return err
	}
	return nil
}

func (logger *FileLogger) closeIfExists(f *os.File) error {
	if f == nil {
		return nil
	}
	return f.Close()
}
