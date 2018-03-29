package core

import (
	"golang.org/x/sync/errgroup"
)

// Commit ...
func (component *Component) Commit(parent *Component) error {
	if len(component.Jobs) == 0 {
		return nil
	}

	if err := component.prefetch(); err != nil {
		return err
	}

	eg := new(errgroup.Group)
	for _, job := range component.Jobs {
		eg.Go(job.Commit)
	}
	return eg.Wait()
}

// prefetch fetches shared data to the shared data instances.
func (component *Component) prefetch() error {

	return nil
}
