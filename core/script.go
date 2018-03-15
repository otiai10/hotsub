package core

// Script specifies what to do inside the container.
// If the Image.Executable is true, whole the "Script" is ignored.
type Script struct {

	// Path is a file path to the script which should be executed on the container.
	Path string

	// Inline is inline command which should be executed on the container.
	// If "Path" is specified, Inline is ignored.
	Inline []string
}
