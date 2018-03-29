package core

import (
	"log"
	"os"
	"time"

	"github.com/otiai10/dkmachine/v0/dkmachine"
	"golang.org/x/sync/errgroup"
)

// Component represents a independent workflow component, handling only 1 input set.
type Component struct {

	// Identity specifies the unique identity of this component.
	Identity Identity

	// Jobs represent specific set of jobs which should be executed on this component.
	Jobs []*Job

	// SharedData ...
	SharedData SharedData

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

	Runtime struct {
		Image  *Image
		Script *Script
	}

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
		Runtime: struct {
			Image  *Image
			Script *Script
		}{Image: &Image{}, Script: &Script{}},
	}
}

// Create ...
func (component *Component) Create() error {

	g := new(errgroup.Group)

	for i, job := range component.Jobs {
		job.Identity.Prefix = component.Identity.Name
		job.Identity.Index = i
		job.Machine.Spec = component.Machine.Spec

		g.Go(job.Create)

		// Delegate runtimes
		job.Container.Image = component.Runtime.Image
		job.Container.Script = component.Runtime.Script
	}

	if len(component.SharedData.Inputs) != 0 {
		g.Go(func() error {
			instance, err := dkmachine.Create(component.SharedData.Spec)
			component.SharedData.Instance = instance
			return err
		})
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

	if component.SharedData.Instance != nil {
		if err := component.SharedData.Instance.Remove(); err != nil {
			e = err
		}
	}

	return e
}
