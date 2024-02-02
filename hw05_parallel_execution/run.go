package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() (string, error)

type counter struct {
	sync.Mutex
	threshold    int
	currentOk    int
	currentError int
	ignoreErrors bool
}

func newCounter(threshold int) *counter {
	return &counter{
		threshold:    threshold,
		ignoreErrors: threshold <= 0,
	}
}

type executionResult string

const (
	resultOk    executionResult = "ok"
	resultError executionResult = "error"
)

func (c *counter) increment(state executionResult) {
	c.Lock()
	defer c.Unlock()

	switch state {
	case resultOk:
		{
			c.currentOk++
		}
	case resultError:
	default:
		{
			c.currentError++
		}
	}
}

func (c *counter) reachedThreshold() bool {
	c.Lock()
	defer c.Unlock()

	return !c.ignoreErrors && c.currentError >= c.threshold
}

func Run(tasks []Task, n int, m int) error {
	ctr := newCounter(m)

	var wg sync.WaitGroup
	wg.Add(n)

	tasksChannel := make(chan Task)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for task := range tasksChannel {
				_, err := task()
				if err != nil {
					ctr.increment(resultError)
				} else {
					ctr.increment(resultOk)
				}
			}
		}()
	}

	for _, task := range tasks {
		if ctr.reachedThreshold() {
			break
		}
		tasksChannel <- task
	}

	close(tasksChannel)

	wg.Wait()

	if ctr.reachedThreshold() {
		return ErrErrorsLimitExceeded
	}
	return nil
}
