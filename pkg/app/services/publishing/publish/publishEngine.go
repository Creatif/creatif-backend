package publish

import (
	"context"
	"fmt"
)

type fnExecutioner struct {
	fn func() error
	// cancellations of queries being executed
	cancel context.CancelFunc
}

type mappedFn = map[string]fnExecutioner

type result struct {
	name  string
	error error
}

type publishEngine struct {
	workers map[string]chan result
}

func (p *publishEngine) addWorker(name string) {
	p.workers[name] = make(chan result)
}

func cancelAll(cancellations []context.CancelFunc) {
	for _, cancel := range cancellations {
		cancel()
	}
}

func (p *publishEngine) run(fns mappedFn, ctx context.Context) map[string]result {
	done := make(chan bool)
	results := make(map[string]result)
	cancellations := make([]context.CancelFunc, 0)
	for _, fn := range fns {
		cancellations = append(cancellations, fn.cancel)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				cancelAll(cancellations)

				return
			}
		}
	}()

	for name, worker := range p.workers {
		go func(name string, worker chan result) {
			results[name] = <-worker

			fmt.Println("Received: ", name)

			if len(results) == 3 {
				done <- true
			}
		}(name, worker)

		go func(name string, worker chan result) {
			fn := fns[name]

			fmt.Println("Sending: ", name)

			err := fn.fn()
			worker <- result{error: err, name: name}
		}(name, worker)
	}

	<-done
	cancelAll(cancellations)
	return results
}

func newPublishEngine() *publishEngine {
	return &publishEngine{
		workers: make(map[string]chan result),
	}
}
