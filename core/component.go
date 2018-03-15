package core

// Component represents a independent workflow component, handling only 1 input set.
type Component struct {

	// Identity specifies the unique identity of this component.
	Identity Identity

	// Jobs represent specific set of jobs which should be executed on this component.
	Jobs []*Job

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
	}
}

// Commit ...
func (component *Component) Commit(parent *Component) error {
	return nil
}
