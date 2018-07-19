package core

// Parameters specifies the parameters assigned to this job.
// It is exactly what the corresponding row in tasks file is parsed to.
type Parameters struct {
	Inputs  Inputs
	Outputs Outputs
	Envs    []Env

	// Includes represent **local** files
	// which should be transfered to VM from local.
	Includes Includes
}
