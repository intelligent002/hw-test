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

const (
	stateOk = iota
	stateError
)

func (c *counter) increment(state int) {
	c.Lock()
	defer c.Unlock()

	switch state {
	case stateOk:
		{
			c.currentOk++
		}
	case stateError:
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
				var _, err = task()
				if err != nil {
					ctr.increment(stateError)
				} else {
					ctr.increment(stateOk)
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
