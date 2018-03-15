package core

// Driver represents Job Driver.
type Driver struct {
	job *Job
}

// Create [lifecycle] create an instance.
func (d *Driver) Create() error {
	return nil
}

// Init [lifecycle] initialize containers inside the instance.
func (d *Driver) Init() error {
	return nil
}

// Attach [lifecycle] mount onto SharedDataInstance.
func (d *Driver) Attach() error {
	return nil
}

// Fetch [lifecycle] download input files and translate URLs to local file paths.
func (d *Driver) Fetch() error {
	return nil
}

// Exec [lifecycle] execute user defined process specified with Image and Script.
func (d *Driver) Exec() error {
	return nil
}

// Push [lifecycle] upload output files to the specified URL.
func (d *Driver) Push() error {
	return nil
}

// Destroy [lifecycle] delete this instance.
func (d *Driver) Destroy() error {
	return nil
}
