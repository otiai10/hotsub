package core

const (
	HOTSUB_HOSTROOT      = "/tmp"
	HOTSUB_CONTAINERROOT = "/tmp"
	HOTSUB_SHARED_DIR    = "__shared"

	// For SharedDataInstances
	HOTSUB_MOUNTPOINT = "/tmp"

	// CreateMaxRetry represents max count for retrying `docker-machine create`
	CreateMaxRetry = 4

	// ContainerMaxRetry represents max count for retrying operations inside docker containers,
	// such as "pull image", "exec create" and "exec create".
	// See https://github.com/otiai10/daap/commit/8b5dfbd93d169c0ae30ce30ea23f81e97f009f7f
	// for more information.
	ContainerMaxRetry = 4
)
