package main

import (
	"net/http"
)

type accountWorkQueueJob struct {
	client              *http.Client
	projectId           string
	accountStructureId  string
	propertyStructureId string
	groupIds            []string
	account             account
}

type accountWorkQueue struct {
	listeners     []chan accountWorkQueueJob
	listWorkQueue propertiesWorkQueue
	jobDoneQueue  chan bool
	balancer      *balancer
}

func newMapWorkQueueJob(
	client *http.Client,
	projectId,
	accountStructureId string,
	propertyStructureId string,
	groupIds []string,
	account account,
) accountWorkQueueJob {
	return accountWorkQueueJob{
		client:              client,
		projectId:           projectId,
		accountStructureId:  accountStructureId,
		propertyStructureId: propertyStructureId,
		groupIds:            groupIds,
		account:             account,
	}
}

func newAccountWorkQueue(workersNum int, buffer int, listWorkQueue propertiesWorkQueue) *accountWorkQueue {
	listeners := make([]chan accountWorkQueueJob, workersNum)
	for i := 0; i < workersNum; i++ {
		listeners[i] = make(chan accountWorkQueueJob, buffer)
	}

	return &accountWorkQueue{
		listeners:     listeners,
		listWorkQueue: listWorkQueue,
		jobDoneQueue:  make(chan bool),
		balancer:      newBalancer(workersNum),
	}
}

func (wq *accountWorkQueue) addJob(j accountWorkQueueJob) {
	worker := wq.balancer.addJob()
	wq.listeners[worker] <- j
}

func (wq *accountWorkQueue) start() chan bool {
	done := make(chan bool)
	for i := 0; i < len(wq.listeners); i++ {
		go func(i int) {
			for {
				select {
				case <-done:
					return
				case j := <-wq.listeners[i]:
					wq.balancer.removeJob(i)
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

							wq.listWorkQueue.addJob(newPropertyWorkQueueJoby(
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
