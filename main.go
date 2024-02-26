package main

import (
	"fmt"
	"sync"
)

type Engine interface {
	Start()
}

type engine[T any] struct {
	jobs     *[][]T
	callback func() error
}

func (e engine[T]) Start() {
	length := 55

	numOfGoroutines := length / 20
	split := length / numOfGoroutines
	remainder := length % numOfGoroutines
	marks := make([][]int, 0)

	for i := 0; i < numOfGoroutines; i++ {
		marks = append(marks, []int{i * split, (i + 1) * split})
		i++
	}

	fmt.Println(marks)

	wg := &sync.WaitGroup{}
	for i := 0; i < len(marks); i++ {
		wg.Add(1)
		go func(marks []int) {

			wg.Done()
		}(marks[i])
	}

	wg.Wait()
}

func NewEngine[T any](jobs *[][]T, cb func() error) Engine {
	return engine[T]{jobs: jobs, callback: cb}
}

func main() {
	jobs := make([][]string, 0)
	currentBatch := make([]string, 0)
	for i := 0; i < 250_567; i++ {
		if len(currentBatch) == 4500 {
			jobs = append(jobs, currentBatch)
			currentBatch = make([]string, 0)
		}

		currentBatch = append(currentBatch, fmt.Sprintf("name-%d", i))
	}

	if len(currentBatch) > 0 {
		jobs = append(jobs, currentBatch)
		currentBatch = nil
	}

	eng := NewEngine[string](&jobs, func() error {
		return nil
	})

	eng.Start()
}
