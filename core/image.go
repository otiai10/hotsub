package core

// Image ...
type Image struct {

	// Name is the container image name with format `IMAGE[:TAG|@DIGEST]`.
	// IMAGE is described as "[registry/]owner/name".
	Name string

	// Executable means if this image has its own ENTRYPOINT or CMD
	// and doesn't need to Exec any other additional script.
	Executable bool
}
