package main

import (
	"math"
	"sync"
)

type balancer struct {
	nodes      map[int]int
	bestWorker int
	lock       sync.Locker
}

func newBalancer(workerNum int) *balancer {
	nodes := make(map[int]int)
	for i := 0; i < workerNum; i++ {
		nodes[i] = 0
	}

	return &balancer{nodes: nodes, lock: &sync.RWMutex{}}
}

func (b *balancer) addJob() int {
	b.lock.Lock()
	defer b.lock.Unlock()
	best := math.MaxUint32
	worker := 0

	for workerId, numOfJobs := range b.nodes {
		if numOfJobs < best {
			best = numOfJobs
			worker = workerId
		}
	}

	b.nodes[worker] = b.nodes[worker] + 1

	return worker
}

func (b *balancer) removeJob(workerId int) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.nodes[workerId] = b.nodes[workerId] - 1
}
