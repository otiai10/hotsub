package core

import (
	"github.com/otiai10/stackerr"
)

// Destroy ...
func (component *Component) Destroy() error {

	errs := stackerr.New()

	if component.SharedData.Instance != nil {
		if err := component.SharedData.Instance.Remove(); err != nil {
			errs.Push(err)
		}
	}

	return errs.Err()
}
