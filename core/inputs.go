package core

// Inputs represent given input set
type Inputs []*Input

// Input represents a input specified as a cell of tasks file.
type Input struct {
	Resource `json:",inline"`
	// // Name is a name label for this input
	// Name string `json:"name"`
	//
	// // Recursive specify if this input is a directory or just a file.
	// Recursive bool `json:"recursive"  yaml:"recursive"`
	//
	// // URL is (in most cases) a resource location through the Internet,
	// // s3://..., gs://..., http://..., https://..., for examples.
	// // The resource specified by this URL would be downloaded to computing node
	// // and translated to local file path.
	// URL string `json:"url"        yaml:"url"`
	// // DeployedPath is a local file path which is translated from URL.
	// DeployedPath string `json:"local_path" yaml:"local_path"`
}
