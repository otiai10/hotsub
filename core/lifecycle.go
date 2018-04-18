package core

// Lifecycle represents the lifecycles of Component and Job.
type Lifecycle string

const (
	// CREATE is creating computing instances, physically.
	CREATE Lifecycle = "CREATE"
	// CONSTRUCT is creating containers inside the instances.
	CONSTRUCT Lifecycle = "CONSTRUCT"
	// FETCH is downloading specified input files from remote storage service.
	FETCH Lifecycle = "FETCH"
	// EXECUTE is executing the specified script inside user-workflow container.
	EXECUTE Lifecycle = "EXECUTE"
	// PUSH is uploading the result files to remote storage service.
	PUSH Lifecycle = "PUSH"
	// DESTROY is deleting the physical instances which are no longer used.
	DESTROY Lifecycle = "DESTROY"
)
