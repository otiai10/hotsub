package core

// Prepare ...
func (component *Component) Prepare() error {
	if len(component.SharedData.Inputs) != 0 {
		if err := component.SharedData.Create(); err != nil {
			return err
		}
	}
	return nil
}
