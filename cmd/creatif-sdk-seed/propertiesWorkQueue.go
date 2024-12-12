package main

import (
	"creatif-sdk-seed/dataGeneration"
	"creatif-sdk-seed/errorHandler"
	"net/http"
)

type propertyWorkQueueJob struct {
	client              *http.Client
	projectId           string
	propertyStructureId string
	propertyVariable    dataGeneration.PropertyVariable
	references          []map[string]string
	imagePaths          []string
}

type propertiesWorkQueue struct {
	listeners    []chan propertyWorkQueueJob
	jobDoneQueue chan bool
	balancer     *balancer
}

func newPropertyWorkQueueJoby(
	client *http.Client,
	projectId,
	propertyStructureId string,
	propertyVariable dataGeneration.PropertyVariable,
	references []map[string]string,
	imagePaths []string,
) propertyWorkQueueJob {
	return propertyWorkQueueJob{
		client:              client,
		projectId:           projectId,
		propertyStructureId: propertyStructureId,
		propertyVariable:    propertyVariable,
		references:          references,
		imagePaths:          imagePaths,
	}
}

func newPropertiesWorkQueue(workersNum int, buffer int) propertiesWorkQueue {
	listeners := make([]chan propertyWorkQueueJob, workersNum)
	for i := 0; i < workersNum; i++ {
		listeners[i] = make(chan propertyWorkQueueJob, buffer)
	}

	return propertiesWorkQueue{
		listeners:    listeners,
		jobDoneQueue: make(chan bool),
		balancer:     newBalancer(workersNum),
	}
}

func (wq propertiesWorkQueue) addJob(j propertyWorkQueueJob) {
	worker := wq.balancer.addJob()
	wq.listeners[worker] <- j
}

func (wq propertiesWorkQueue) start() chan bool {
	done := make(chan bool)
	for i := 0; i < len(wq.listeners); i++ {
		go func(i int) {
			for {
				select {
				case <-done:
					return
				case j := <-wq.listeners[i]:
					errorHandler.HandleHttpError(createProperty(
						j.client,
						j.projectId,
						j.propertyStructureId,
						j.propertyVariable,
						j.references,
						j.imagePaths,
					))

					wq.balancer.removeJob(i)

					wq.jobDoneQueue <- true
				}
			}
		}(i)

	}

	return done
}
