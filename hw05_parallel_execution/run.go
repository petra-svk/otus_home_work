package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(wg *sync.WaitGroup, queue <-chan Task, result chan<- error) {
	defer wg.Done()

	for task := range queue {
		result <- task()
	}
}

func Run(tasks []Task, n, m int) error {
	wg := new(sync.WaitGroup)
	queue := make(chan Task)   // queue with tasks
	result := make(chan error) // queue where workers put result of tasks
	var countError int

	if len(tasks) < n {
		n = len(tasks)
	}
	// set up workers
	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(wg, queue, result)
	}

	countTask := 0

LOOP:
	j := 0
	for i := 0; i < n; i++ {
		queue <- tasks[countTask]
		countTask++
		j++
		if countTask == len(tasks) {
			break
		}
	}

	for i := 0; i < j; i++ {
		res := <-result
		if m > 0 && res != nil {
			countError++
		}
	}

	if (m > 0 && countError >= m) || countTask == len(tasks) {
		close(queue)
	} else {
		goto LOOP
	}

	wg.Wait()
	close(result)

	if m > 0 && countError >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
