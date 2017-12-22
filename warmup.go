package main

import (
	"context"
	"sync"

	"github.com/otiai10/daap"
)

// Warmup ...
func (h *Handler) Warmup(ctx context.Context, c *daap.Container, job *Job, wg *sync.WaitGroup) {
	defer wg.Done()
	job.Logf("Pulling the image to the host: %v", c.Image)
	pull, err := c.PullImage(ctx)
	if err != nil {
		job.Errorf("failed to pull image: %v", err)
		return
	}
	for range pull {
		// fmt.Print(".")
	}
	job.Logf("The image pulled successfully: %v", c.Image)

	job.Logf("Creating a container on the host")
	if err = c.Create(ctx); err != nil {
		job.Errorf("failed to create container: %v: %v", c.Image, err)
		return
	}
	job.Logf("The container created successfully")

	job.Logf("Starting the container up")
	if err = c.Start(ctx); err != nil {
		job.Errorf("failed to start container: %v: %v", c.Image, err)
		return
	}
	job.Logf("The container started successfully: %v", c.Image)
}
