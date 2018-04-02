package core

import (
	"golang.org/x/sync/errgroup"
)

// Commit ...
func (component *Component) Commit(parent *Component) error {
	if len(component.Jobs) == 0 {
		return nil
	}

	eg := new(errgroup.Group)
	for _, job := range component.Jobs {
		eg.Go(job.Commit)
	}
	return eg.Wait()
}
