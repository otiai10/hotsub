package core

import (
	"log"
	"os"
	"time"
)

// Component represents a independent workflow component, handling only 1 input set.
type Component struct {

	// Identity specifies the unique identity of this component.
	Identity Identity

	// Jobs represent specific set of jobs which should be executed on this component.
	Jobs []*Job

	// SharedData ...
	SharedData *SharedData

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
		SharedData: &SharedData{},
	}
}
