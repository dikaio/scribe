package build

import (
	"runtime"
	"sync"
)

// parallelWorker represents a generic worker function
type parallelWorker func(workerID int, jobs <-chan interface{}, results chan<- interface{}, errChan chan<- error, wg *sync.WaitGroup)

// parallelExecutor runs tasks in parallel using a worker pool
func parallelExecutor(jobs []interface{}, worker parallelWorker) ([]interface{}, []error) {
	numWorkers := runtime.NumCPU()
	if len(jobs) < numWorkers {
		numWorkers = len(jobs)
	}

	// Create channels
	jobChan := make(chan interface{}, len(jobs))
	resultsChan := make(chan interface{}, len(jobs))
	errChan := make(chan error, len(jobs))

	// Add jobs to channel
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan)

	// Start workers
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(i, jobChan, resultsChan, errChan, &wg)
	}

	// Wait for workers to finish
	wg.Wait()
	close(resultsChan)
	close(errChan)

	// Collect results and errors
	results := make([]interface{}, 0, len(jobs))
	for result := range resultsChan {
		results = append(results, result)
	}

	errors := make([]error, 0)
	for err := range errChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	return results, errors
}