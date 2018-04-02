package core

// Destroy ...
func (component *Component) Destroy() error {

	var e error

	for _, job := range component.Jobs {
		if err := job.Destroy(); err != nil {
			e = err
		}
	}

	if component.SharedData.Instance != nil {
		if err := component.SharedData.Instance.Remove(); err != nil {
			e = err
		}
	}

	return e
}
