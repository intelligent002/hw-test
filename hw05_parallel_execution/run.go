package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() (string, error)

type counter struct {
	sync.Mutex
	threshold    int
	current      int
	ignoreErrors bool
}

func newCounter(threshold int) *counter {
	return &counter{
		threshold:    threshold,
		ignoreErrors: threshold <= 0,
	}
}

func (c *counter) increment() {
	c.Lock()
	defer c.Unlock()

	c.current++
}

func (c *counter) reachedThreshold() bool {
	c.Lock()
	defer c.Unlock()

	return !c.ignoreErrors && c.current >= c.threshold
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
				var res, err = task()
				if err != nil {
					fmt.Println(err)
					ctr.increment()
				} else {
					fmt.Println(res)
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
