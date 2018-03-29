package core

// Destroy ...
func (job *Job) Destroy() error {
	if job.Machine.Instance == nil {
		return nil
	}
	return job.Machine.Instance.Remove()
}
