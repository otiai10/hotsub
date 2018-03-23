package core

import (
	"log"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

// Component represents a independent workflow component, handling only 1 input set.
type Component struct {

	// Identity specifies the unique identity of this component.
	Identity Identity

	// Jobs represent specific set of jobs which should be executed on this component.
	Jobs []*Job

	// Machine represents the spec of machines on which each job is executed.
	Machine *Machine

	/* You Ain't Gonna Need It!! */
	// // Nest can specify nested components.
	// // If "Nest" is provided not nil, all the "Jobs" are ignored.
	// // If neither "Parallel" nor "Serial" is provided, it results in an error.
	// Nest *struct {
	// 	Parallel []*Component
	// 	Serial   []*Component
	// }

	/* You Ain't Gonna Need It!! */
	// // Parent
	// Parent *Component

	// Report directory path
	Report struct {
		// LocalPath is a local path to save report files.
		LocalPath string
		// URL, if specified, the report path would be uploaded to this URL.
		URL string
		// Message is an interface to write log
	}

	// Log is an application logger ONLY FOR ROOT COMPONENT.
	Log *log.Logger
}

// RootComponentTemplate ...
func RootComponentTemplate(name string) *Component {
	return &Component{
		Identity: Identity{Name: name, Timestamp: time.Now().UnixNano()},
		Log:      log.New(os.Stdout, "[root]", 1),
		Machine:  &Machine{},
	}
}

// Commit ...
func (component *Component) Commit(parent *Component) error {
	if len(component.Jobs) == 0 {
		return nil
	}

	if err := component.Create(); err != nil {
		return err
	}
	defer component.Destroy()

	return nil
}

// Create ...
func (component *Component) Create() error {
	g := new(errgroup.Group)
	for i, job := range component.Jobs {
		job.Identity.Prefix = component.Identity.Name
		job.Identity.Index = i
		job.Machine.Spec = component.Machine.Spec
		g.Go(job.Create)
	}
	return g.Wait()
}

// Destroy ...
func (component *Component) Destroy() error {
	var e error
	for _, job := range component.Jobs {
		if err := job.Destroy(); err != nil {
			e = err
		}
	}
	return e
}
