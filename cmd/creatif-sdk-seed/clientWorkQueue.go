package main

import (
	"creatif-sdk-seed/dataGeneration"
	"creatif-sdk-seed/errorHandler"
	"net/http"
)

type clientWorkQueueJob struct {
	client              *http.Client
	projectId           string
	clientStructureId   string
	propertyStructureId string
	groupIds            []string
	clientVariable      dataGeneration.Client
}

type clientWorkQueue struct {
	listeners           []chan clientWorkQueueJob
	listWorkQueue       propertiesWorkQueue
	jobDoneQueue        chan bool
	balancer            *balancer
	propertiesPerStatus int
}

func newClientWorkQueueJob(
	client *http.Client,
	projectId,
	clientStructureId string,
	propertyStructureId string,
	groupIds []string,
	clientVariable dataGeneration.Client,
) clientWorkQueueJob {
	return clientWorkQueueJob{
		client:              client,
		projectId:           projectId,
		clientStructureId:   clientStructureId,
		propertyStructureId: propertyStructureId,
		groupIds:            groupIds,
		clientVariable:      clientVariable,
	}
}

func newClientWorkQueue(workersNum int, buffer int, listWorkQueue propertiesWorkQueue, propertiesPerStatus int) *clientWorkQueue {
	listeners := make([]chan clientWorkQueueJob, workersNum)
	for i := 0; i < workersNum; i++ {
		listeners[i] = make(chan clientWorkQueueJob, buffer)
	}

	return &clientWorkQueue{
		listeners:           listeners,
		listWorkQueue:       listWorkQueue,
		jobDoneQueue:        make(chan bool),
		balancer:            newBalancer(workersNum),
		propertiesPerStatus: propertiesPerStatus,
	}
}

func (wq *clientWorkQueue) addJob(j clientWorkQueueJob) {
	worker := wq.balancer.addJob()
	wq.listeners[worker] <- j
}

func (wq *clientWorkQueue) start() chan bool {
	done := make(chan bool)
	for i := 0; i < len(wq.listeners); i++ {
		go func(i int) {
			for {
				select {
				case <-done:
					return
				case j := <-wq.listeners[i]:
					clientId := addToMapAndGetClientId(
						j.client,
						j.projectId,
						j.clientStructureId,
						j.clientVariable,
					)

					wq.jobDoneQueue <- true

					wq.balancer.removeJob(i)

					propertiesGen := dataGeneration.NewPropertiesGenerator()
					for {
						newSequence, ok := propertiesGen.Generate()

						if !ok {
							break
						}

						/**
						There are 5 locales, 3 property statuses and 4 property types. For every property type, generate 10 properties.
						Every locale generate 120 properties. That is property status * property type * num of locales * 10.
						The calculation is then 5 * 3 * 4 * 10
						*/
						for a := 0; a < wq.propertiesPerStatus; a++ {
							singleProperty, err := dataGeneration.GenerateSingleProperty(clientId, newSequence.Locale, newSequence.PropertyStatus, newSequence.PropertyType, j.groupIds)
							if err != nil {
								errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
							}

							wq.listWorkQueue.addJob(newPropertyWorkQueueJoby(
								j.client,
								j.projectId,
								j.propertyStructureId,
								singleProperty.Variable,
								singleProperty.Connections,
								singleProperty.ImagePaths,
							))
						}
					}
				}
			}
		}(i)
	}

	return done
}
