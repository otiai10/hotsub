// package main
// But, this module should be in a separated package in the future.
package main

func handleBunch(tasks []*Task) <-chan Job {

	results := make(chan Job)
	done := make(chan bool)
	count := 0

	for _, task := range tasks {
		go func(t *Task) {
			results <- handle(t)
			done <- true
		}(task)
	}

	go func() {
		defer close(done)
		defer close(results)
		for {
			<-done
			count++
			if count >= len(tasks) {
				return
			}
		}
	}()

	return results
}

func handle(task *Task) Job {
	job := Job{Task: *task}

	// {{{ TODO: Refactor and Slim up

	// }}}

	return job
}
