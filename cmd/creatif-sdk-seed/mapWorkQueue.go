package main

import (
	"net/http"
)

type mapWorkQueueJob struct {
	client              *http.Client
	projectId           string
	accountStructureId  string
	propertyStructureId string
	groupIds            []string
	account             account
}

type mapWorkQueue struct {
	listeners     []chan mapWorkQueueJob
	listWorkQueue listWorkQueue
	jobDoneQueue  chan bool
}

func newMapWorkQueueJob(
	client *http.Client,
	projectId,
	accountStructureId string,
	propertyStructureId string,
	groupIds []string,
	account account,
) mapWorkQueueJob {
	return mapWorkQueueJob{
		client:              client,
		projectId:           projectId,
		accountStructureId:  accountStructureId,
		propertyStructureId: propertyStructureId,
		groupIds:            groupIds,
		account:             account,
	}
}

func newMapWorkQueue(workersNum int, buffer int, listWorkQueue listWorkQueue) *mapWorkQueue {
	listeners := make([]chan mapWorkQueueJob, workersNum)
	for i := 0; i < workersNum; i++ {
		listeners[i] = make(chan mapWorkQueueJob, buffer)
	}

	return &mapWorkQueue{
		listeners:     listeners,
		listWorkQueue: listWorkQueue,
		jobDoneQueue:  make(chan bool),
	}
}

func (wq *mapWorkQueue) addJob(j mapWorkQueueJob) {
	worker := randomBetween(1, len(wq.listeners)-1)
	wq.listeners[worker] <- j
}

func (wq *mapWorkQueue) start() chan bool {
	done := make(chan bool)
	for i := 0; i < len(wq.listeners); i++ {
		go func(i int) {
			for {
				select {
				case <-done:
					break
				case j := <-wq.listeners[i]:
					accountId := addToMapAndGetAccountId(
						j.client,
						j.projectId,
						j.accountStructureId,
						j.account,
					)

					wq.jobDoneQueue <- true

					propertiesGen := newPropertiesGenerator()
					for {
						newSequence, ok := propertiesGen.generate()
						if !ok {
							break
						}
						
						for a := 0; a < 10; a++ {
							singleProperty, err := generateSingleProperty(accountId, newSequence.locale, newSequence.propertyStatus, newSequence.propertyType, j.groupIds)
							if err != nil {
								handleAppError(err, Cannot_Continue_Procedure)
							}

							wq.listWorkQueue.addJob(newListWorkQueueJob(
								j.client,
								j.projectId,
								j.propertyStructureId,
								singleProperty.variable,
								singleProperty.references,
								singleProperty.imagePaths,
							))
						}
					}
				}
			}
		}(i)
	}

	return done
}
