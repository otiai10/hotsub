package platform

import (
	"github.com/otiai10/hotsub/params"
)

// Virtualbox ...
type Virtualbox struct {
}

// Validate does nothing so far.
// If needed, create a network interface or something like that
// in the future.
func (p *Virtualbox) Validate(ctx params.Context) error {
	return nil
}

// HyperV ...
type HyperV struct {
}

// Validate does nothing so far.
// If needed, create a network interface or something like that
// in the future.
func (p *HyperV) Validate(ctx params.Context) error {
	return nil
}
