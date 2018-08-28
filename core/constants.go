package core

const (

	// HotsubHostRoot is a mount point mounted by containers on this host machine.
	HotsubHostRoot = "/tmp"

	// HotsubContainerRoot is a root directory inside each container.
	// This path is passed as "HOTSUB_ROOT" to workflows,
	// wish you don't need to refer "HOTSUB_ROOT" inside your workflow.
	HotsubContainerRoot = "/tmp"

	// HotsubSharedDirectoryPath is a path-prefix from ContainerRoot,
	// in which all the shared data are located.
	HotsubSharedDirectoryPath = "__shared"

	// HotsubSharedInstanceMountPoint is a mount point mountedd by containers of computing instances.
	HotsubSharedInstanceMountPoint = "/tmp"

	// CreateMaxRetry represents max count for retrying `docker-machine create`
	CreateMaxRetry = 4

	// ContainerMaxRetry represents max count for retrying operations inside docker containers,
	// such as "pull image", "exec create" and "exec create".
	// See https://github.com/otiai10/daap/commit/8b5dfbd93d169c0ae30ce30ea23f81e97f009f7f
	// for more information.
	ContainerMaxRetry = 4
)
