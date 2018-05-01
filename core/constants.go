package core

const (
	AWSUB_HOSTROOT      = "/tmp"
	AWSUB_CONTAINERROOT = "/tmp"
	AWSUB_SHARED_DIR    = "__shared"

	// For SharedDataInstances
	AWSUB_MOUNTPOINT = "/tmp"

	// ContainerMaxRetry represents max count for retrying operations inside docker containers,
	// such as "pull image", "exec create" and "exec create".
	// See https://github.com/otiai10/daap/commit/8b5dfbd93d169c0ae30ce30ea23f81e97f009f7f
	// for more information.
	ContainerMaxRetry = 4
)
