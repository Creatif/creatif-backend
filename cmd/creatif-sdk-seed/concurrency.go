package main

import (
	"creatif-sdk-seed/dataGeneration"
	"creatif-sdk-seed/errorHandler"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type joinedStructureClient struct {
	projectId           string
	groupIds            []string
	clientId            string
	clientStructureId   string
	propertyStructureId string

	client dataGeneration.Client
}

type projectProduct struct {
	projectId           string
	groupIds            []string
	clientStructureId   string
	propertyStructureId string
}

func newJoinedStructureClient(
	projectId,
	clientStructureId,
	propertyStructureId string,
	groupIds []string,

	client dataGeneration.Client,
) joinedStructureClient {
	return joinedStructureClient{
		projectId:           projectId,
		groupIds:            groupIds,
		clientStructureId:   clientStructureId,
		propertyStructureId: propertyStructureId,
		client:              client,
	}
}

func newProjectProduct(projectId, clientStructureId, propertyStructureId string, groupIds []string) projectProduct {
	return projectProduct{
		projectId:           projectId,
		groupIds:            groupIds,
		clientStructureId:   clientStructureId,
		propertyStructureId: propertyStructureId,
	}
}

func projectProducer(client *http.Client, numOfProjects int) []chan projectProduct {
	producers := make([]chan projectProduct, numOfProjects)
	for i := 0; i < numOfProjects; i++ {
		producers[i] = make(chan projectProduct)
	}

	go func() {
		projectNames := []string{
			"Warsaw Brokers",
			"London Brokers",
			"Paris Brokers",
			"Berlin Brokers",
			"Barcelona Brokers",
			"Zagreb Brokers",
			"Belgrade Brokers",
			"Prag Brokers",
			"Rome Brokers",
			"Athens Brokers",
		}

		for i := 0; i < numOfProjects; i++ {
			projectName := projectNames[i]
			var projectId string
			result := errorHandler.HandleHttpError(createProject(client, projectName))
			res := result.Response()

			if res != nil && res.Body == nil {
				errorHandler.HandleAppError(errors.New("projectProducer() is trying to work on a nil body"), Cannot_Continue_Procedure)
			}

			defer res.Body.Close()
			var m map[string]interface{}
			b, err := io.ReadAll(res.Body)
			if err != nil {
				errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
			}

			if err := json.Unmarshal(b, &m); err != nil {
				errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
			}

			if res.StatusCode < 200 && res.StatusCode > 299 {
				errorHandler.HandleAppError(errors.New(fmt.Sprintf("Creating project failed with status %d and body %s", res.StatusCode, string(b))), Cannot_Continue_Procedure)
			}

			projectId = m["id"].(string)

			groupIds := createGroupsAndGetGroupIds(client, projectId)
			clientStructureId := createClientStructureAndReturnId(client, projectId)
			propertyStructureId := createPropertiesStructureAndReturnID(client, projectId)

			producers[i] <- newProjectProduct(projectId, clientStructureId, propertyStructureId, groupIds)
			close(producers[i])
		}
	}()

	return producers
}

func clientProducer(client *http.Client, projectProducers []chan projectProduct, wq *clientWorkQueue, reporter *reporter) []projectProduct {
	publishingListeners := make([]projectProduct, len(projectProducers))
	for i, producer := range projectProducers {
		projectProductResult := <-producer
		publishingListeners[i] = projectProductResult

		reporter.AddProjectID(projectProductResult.projectId)

		groupIds := projectProductResult.groupIds
		projectId := projectProductResult.projectId
		clientStructureId := projectProductResult.clientStructureId
		propertyStructureId := projectProductResult.propertyStructureId

		for a := 0; a < 100; a++ {
			getClient, err := dataGeneration.GenerateSingleClient(groupIds)
			if err != nil {
				errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
			}

			joinedClient := newJoinedStructureClient(
				projectId,
				clientStructureId,
				propertyStructureId,
				groupIds,
				getClient,
			)

			wq.addJob(newClientWorkQueueJob(
				client,
				joinedClient.projectId,
				joinedClient.clientStructureId,
				joinedClient.propertyStructureId,
				joinedClient.groupIds,
				joinedClient.client,
			))
		}
	}

	return publishingListeners
}

func concurrencyCoordinator(
	propertiesWorkQueue propertiesWorkQueue,
	clientWorkQueue *clientWorkQueue,
	progressBarNotifier chan bool,
	propertyWorkQueueDone chan bool,
	clientWorkQueueDone chan bool,
	reporter *reporter,
) {
	go func() {
		workQueueTimeout := time.After(2 * time.Second)
		for {
			select {
			case <-propertiesWorkQueue.jobDoneQueue:
				progressBarNotifier <- true
				reporter.AddProperty()
				workQueueTimeout = time.After(5 * time.Second)
			case <-clientWorkQueue.jobDoneQueue:
				progressBarNotifier <- true
				reporter.AddClient()
				workQueueTimeout = time.After(5 * time.Second)
			case <-workQueueTimeout:
				close(propertyWorkQueueDone)
				close(clientWorkQueueDone)
				return
			}
		}
	}()
}

func mergeDoneQueues(clientQueue chan bool, propertyQueue chan bool) chan bool {
	done := make(chan bool)

	go func() {
		<-clientQueue
		<-propertyQueue

		done <- true
	}()

	return done
}
