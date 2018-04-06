package core

import (
	"golang.org/x/sync/errgroup"
)

// Construct sets up all the containers inside computing nodes.
func (component *Component) Construct() error {
	eg := new(errgroup.Group)

	for _, job := range component.Jobs {
		j := job
		eg.Go(func() error {
			return j.Construct(component.SharedData)
		})
	}

	return eg.Wait()
}
