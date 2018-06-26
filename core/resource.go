package core

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/otiai10/ternary"
)

// Resource represents a common struct for both Input and Output.
type Resource struct {
	// Name is a name label for this output
	Name string `json:"name"`

	// Recursive specify if this input is a directory or just a file.
	Recursive bool `json:"recursive"  yaml:"recursive"`

	// URL is (in most cases) a resource location through the Internet,
	// s3://..., gs://... for examples.
	// The output location specified by this URL would be translated to
	// local file path on computing node, and pushed to this URL after the job.
	URL string `json:"url"        yaml:"url"`
	// LocalPath is a local file path which is translated from URL.
	LocalPath string `json:"local_path" yaml:"local_path"`
}

// Localize convert given resource URL to local file path inside the container.
func (resource *Resource) Localize(rootdir string) error {
	u, err := url.Parse(resource.URL)
	if err != nil {
		return err
	}
	bucket := u.Host
	resource.LocalPath = filepath.ToSlash(filepath.Join(rootdir, bucket, u.Path))
	return nil
}

// Env ...
func (resource *Resource) Env() Env {
	return Env{
		Name:  resource.Name,
		Value: resource.LocalPath,
	}
}

// EnvForFetch ...
func (resource *Resource) EnvForFetch() []string {
	key := ternary.If(resource.Recursive).String("INPUT_RECURSIVE", "INPUT")
	return []string{
		fmt.Sprintf("%s=%s", key, resource.URL),
		fmt.Sprintf("%s=%s", "DIR", filepath.Dir(resource.LocalPath)),
	}
}
