package core

import (
	"fmt"

	"golang.org/x/sync/errgroup"
)

// Attach attaches SharedDataInstances to computing instances.
func (component *Component) Attach() error {

	shared := component.SharedData

	if len(shared.Inputs) == 0 {
		return nil // Do nothing
	}

	if shared.Instance == nil {
		return fmt.Errorf("couldn't get shared data instances")
	}

	env := []Env{}
	for _, input := range shared.Inputs {
		if err := input.Localize(AWSUB_CONTAINERROOT + AWSUB_SHARED_DIR); err != nil {
			return err
		}
		env = append(env, input.Env())
	}

	eg := new(errgroup.Group)
	for _, job := range component.Jobs {
		j := job
		j.Parameters.Envs = append(j.Parameters.Envs, env...)
		eg.Go(func() error { return j.Mount(shared.Instance) })
	}

	return eg.Wait()
}
