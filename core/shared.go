package core

import "github.com/otiai10/dkmachine/v0/dkmachine"

// SharedData ...
type SharedData struct {
	Spec     *dkmachine.CreateOptions
	Instance *dkmachine.Machine
	Inputs   Inputs
}
