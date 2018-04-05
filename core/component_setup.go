package core

import (
	"golang.org/x/sync/errgroup"
)

// Setup sets up all the containers inside computing nodes.
func (component *Component) Setup() error {
	eg := new(errgroup.Group)

	for _, job := range component.Jobs {
		j := job
		eg.Go(func() error {
			return j.Setup(component.SharedData)
		})
	}

	return eg.Wait()
}
