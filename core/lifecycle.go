package core

// Lifecycle represents the lifecycles of Component and Job.
type Lifecycle string

const (
	// CREATE is creating computing instances, physically.
	CREATE Lifecycle = "create"
	// CONSTRUCT is creating containers inside the instances.
	CONSTRUCT Lifecycle = "construct"
	// FETCH is downloading specified input files from remote storage service.
	FETCH Lifecycle = "fetch"
	// EXECUTE is executing the specified script inside user-workflow container.
	EXECUTE Lifecycle = "execute"
	// PUSH is uploading the result files to remote storage service.
	PUSH Lifecycle = "push"
	// DESTROY is deleting the physical instances which are no longer used.
	DESTROY Lifecycle = "destroy"
)
