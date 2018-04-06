package core

// Destroy ...
func (job *Job) Destroy() error {
	job.Lifetimef("[DESTROY]\tTerminating computing instance for this job...")
	if job.Machine.Instance == nil {
		return nil
	}
	return job.Machine.Instance.Remove()
}
