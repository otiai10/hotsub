package core

// Includes represent local files which should be transfered to VM.
type Includes []*Include

// Include represents a file which should be included and transferd to VM.
type Include struct {

	// Because `Include` model represents client local path,
	// `URL` of `Include` does NOT have any meaning.
	Resource `json:",inline"`

	// LocalPath represents a file path in the client machine,
	// where `hotsub` command is issued.
	// This file is `docker cp` to VM and traslated to `DeployedPath`.
	LocalPath string
}
