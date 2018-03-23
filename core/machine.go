package core

import "github.com/otiai10/dkmachine/v0/dkmachine"

// Machine ...
type Machine struct {

	// {{{ TODO: Not used so far
	// Provider specifies machine provider, either of [ec2] ~~[gce, k8s, local]~~
	Provider string
	// CPU specifies how many CPU cores are required for the "Job"
	CPU int
	// Memory specifies how much memory are required (in GB) for the "Job"
	Memory string
	// }}}

	// Spec represent options to create instance.
	Spec *dkmachine.CreateOptions
}

// Instantiate ...
func (m *Machine) Instantiate() (*dkmachine.Machine, error) {
	return dkmachine.Create(m.Spec)
}
