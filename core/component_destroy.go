package core

// Destroy ...
func (component *Component) Destroy() error {

	if component.SharedData.Instance != nil {
		if err := component.SharedData.Instance.Remove(); err != nil {
			return err
		}
	}

	return nil
}
