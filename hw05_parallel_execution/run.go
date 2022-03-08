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
	doneChan := make(chan struct{})
	var errorsCount int32

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			worker(doneChan, tasksChan, &errorsCount)
		}()
	}

	for _, task := range tasks {
		if m > 0 && int(atomic.LoadInt32(&errorsCount)) >= m {
			returnErr = ErrErrorsLimitExceeded
			break
		}

		tasksChan <- task
	}

	// если получили ошибку в последней задаче
	if m > 0 && int(atomic.LoadInt32(&errorsCount)) >= m {
		returnErr = ErrErrorsLimitExceeded
	}

	close(doneChan)
	close(tasksChan)
	wg.Wait()

	return returnErr
}

func worker(doneChan <-chan struct{}, tasksChan <-chan Task, errorsCount *int32) {
	for {
		select {
		case <-doneChan:
			return
		default:
		}

		select {
		case <-doneChan:
			return
		case task, ok := <-tasksChan:
			if !ok {
				return
			}

			if err := task(); err != nil {
				atomic.AddInt32(errorsCount, 1)
			}
		}
	}
}
