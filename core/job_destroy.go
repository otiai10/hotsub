package core

// Destroy ...
func (job *Job) Destroy() error {

	job.Lifetime(DESTROY, "Terminating computing instance for this job...")

	if job.Machine.Instance == nil {
		return nil
	}

	if err := job.Machine.Instance.Remove(); err != nil {
		return err
	}

	if job.Report.Log == nil {
		return nil
	}

	if err := job.Report.Log.Close(); err != nil {
		return err
	}

	return nil
}
