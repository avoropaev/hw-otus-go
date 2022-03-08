package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) (returnErr error) {
	tasksChan := make(chan Task)
	var errorsCount int32

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for task := range tasksChan {
				if err := task(); err != nil && m > 0 {
					atomic.AddInt32(&errorsCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		tasksChan <- task

		if m > 0 && int(atomic.LoadInt32(&errorsCount)) >= m {
			returnErr = ErrErrorsLimitExceeded
			break
		}
	}

	close(tasksChan)
	wg.Wait()

	return returnErr
}
