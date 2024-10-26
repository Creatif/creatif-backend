package main

import (
	"net/http"
)

type listWorkQueueJob struct {
	client              *http.Client
	projectId           string
	propertyStructureId string
	propertyVariable    propertyVariable
	references          []map[string]string
	imagePaths          []string
}

type listWorkQueue struct {
	listeners    []chan listWorkQueueJob
	jobDoneQueue chan bool
}

func newListWorkQueueJob(
	client *http.Client,
	projectId,
	propertyStructureId string,
	propertyVariable propertyVariable,
	references []map[string]string,
	imagePaths []string,
) listWorkQueueJob {
	return listWorkQueueJob{
		client:              client,
		projectId:           projectId,
		propertyStructureId: propertyStructureId,
		propertyVariable:    propertyVariable,
		references:          references,
		imagePaths:          imagePaths,
	}
}

func newListWorkQueue(workersNum int, buffer int) listWorkQueue {
	listeners := make([]chan listWorkQueueJob, workersNum)
	for i := 0; i < workersNum; i++ {
		listeners[i] = make(chan listWorkQueueJob, buffer)
	}

	return listWorkQueue{
		listeners:    listeners,
		jobDoneQueue: make(chan bool),
	}
}

func (wq listWorkQueue) addJob(j listWorkQueueJob) {
	worker := randomBetween(1, len(wq.listeners)-1)
	wq.listeners[worker] <- j
}

func (wq listWorkQueue) start() chan bool {
	done := make(chan bool)
	for i := 0; i < len(wq.listeners); i++ {
		go func(i int) {
			for {
				select {
				case <-done:
					break
				case j := <-wq.listeners[i]:
					handleHttpError(addToList(
						j.client,
						j.projectId,
						j.propertyStructureId,
						j.propertyVariable,
						j.references,
						j.imagePaths,
					))

					wq.jobDoneQueue <- true
				}
			}
		}(i)

	}

	return done
}
