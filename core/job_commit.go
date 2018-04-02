package core

// Commit represents a main process of this job.
// The main process of this job consists of Fetch, Exec, and Push.
func (job *Job) Commit() error {

	if err := job.Fetch(); err != nil {
		return err
	}

	if err := job.Exec(); err != nil {
		return err
	}

	if err := job.Push(); err != nil {
		return err
	}

	return nil
}
