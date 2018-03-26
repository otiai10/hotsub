package core

// Outputs represent given input set
type Outputs []*Output

// Output represents a input specified as a cell of tasks file.
type Output struct {
	Resource `json:",inline"`
	// // Name is a name label for this output
	// Name string `json:"name"`
	//
	// // Recursive specify if this input is a directory or just a file.
	// Recursive bool `json:"recursive"  yaml:"recursive"`
	//
	// // URL is (in most cases) a resource location through the Internet,
	// // s3://..., gs://... for examples.
	// // The output location specified by this URL would be translated to
	// // local file path on computing node, and pushed to this URL after the job.
	// URL string `json:"url"        yaml:"url"`
	// // LocalPath is a local file path which is translated from URL.
	// LocalPath string `json:"local_path" yaml:"local_path"`
}
