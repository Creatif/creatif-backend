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
	listeners           []chan accountWorkQueueJob
	listWorkQueue       propertiesWorkQueue
	jobDoneQueue        chan bool
	balancer            *balancer
	propertiesPerStatus int
}

func newAccountWorkQueueJob(
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

func newAccountWorkQueue(workersNum int, buffer int, listWorkQueue propertiesWorkQueue, propertiesPerStatus int) *accountWorkQueue {
	listeners := make([]chan accountWorkQueueJob, workersNum)
	for i := 0; i < workersNum; i++ {
		listeners[i] = make(chan accountWorkQueueJob, buffer)
	}

	return &accountWorkQueue{
		listeners:           listeners,
		listWorkQueue:       listWorkQueue,
		jobDoneQueue:        make(chan bool),
		balancer:            newBalancer(workersNum),
		propertiesPerStatus: propertiesPerStatus,
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
					accountId := addToMapAndGetAccountId(
						j.client,
						j.projectId,
						j.accountStructureId,
						j.account,
					)

					wq.jobDoneQueue <- true

					wq.balancer.removeJob(i)

					propertiesGen := newPropertiesGenerator()
					for {
						newSequence, ok := propertiesGen.generate()

						if !ok {
							break
						}

						/**
						There are 5 locales, 3 property statuses and 4 property types. For every property type, generate 10 properties.
						Every locale generate 120 properties. That is property status * property type * num of locales * 10.
						The calculation is then 5 * 3 * 4 * 10
						*/
						for a := 0; a < wq.propertiesPerStatus; a++ {
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
