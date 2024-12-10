package main

import (
	"creatif-sdk-seed/dataGeneration"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type joinedStructureAccount struct {
	projectId           string
	groupIds            []string
	accountId           string
	accountStructureId  string
	propertyStructureId string

	account dataGeneration.Account
}

type projectProduct struct {
	projectId           string
	groupIds            []string
	accountStructureId  string
	propertyStructureId string
}

type accountProduct struct {
	projectId           string
	groupIds            []string
	accountId           string
	propertyStructureId string
}

func newJoinedStructureAccount(
	projectId,
	accountStructureId,
	propertyStructureId string,
	groupIds []string,

	account dataGeneration.Account,
) joinedStructureAccount {
	return joinedStructureAccount{
		projectId:           projectId,
		groupIds:            groupIds,
		accountStructureId:  accountStructureId,
		propertyStructureId: propertyStructureId,
		account:             account,
	}
}

func newProjectProduct(projectId, accountStructureId, propertyStructureId string, groupIds []string) projectProduct {
	return projectProduct{
		projectId:           projectId,
		groupIds:            groupIds,
		accountStructureId:  accountStructureId,
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
			result := handleHttpError(createProject(client, projectName))
			res := result.Response()

			if res != nil && res.Body == nil {
				handleAppError(errors.New("projectProducer() is trying to work on a nil body"), Cannot_Continue_Procedure)
			}

			defer res.Body.Close()
			var m map[string]interface{}
			b, err := io.ReadAll(res.Body)
			if err != nil {
				handleAppError(err, Cannot_Continue_Procedure)
			}

			if err := json.Unmarshal(b, &m); err != nil {
				handleAppError(err, Cannot_Continue_Procedure)
			}

			if res.StatusCode < 200 && res.StatusCode > 299 {
				handleAppError(errors.New(fmt.Sprintf("Creating project failed with status %d and body %s", res.StatusCode, string(b))), Cannot_Continue_Procedure)
			}

			projectId = m["id"].(string)

			groupIds := createGroupsAndGetGroupIds(client, projectId)
			accountStructureId := createAccountStructureAndReturnID(client, projectId)
			propertyStructureId := createPropertiesStructureAndReturnID(client, projectId)

			producers[i] <- newProjectProduct(projectId, accountStructureId, propertyStructureId, groupIds)
			close(producers[i])
		}
	}()

	return producers
}

func accountProducer(client *http.Client, projectProducers []chan projectProduct, wq *accountWorkQueue, reporter *reporter) []projectProduct {
	publishingListeners := make([]projectProduct, len(projectProducers))
	for i, producer := range projectProducers {
		projectProductResult := <-producer
		publishingListeners[i] = projectProductResult

		reporter.AddProjectID(projectProductResult.projectId)

		groupIds := projectProductResult.groupIds
		projectId := projectProductResult.projectId
		accountStructureId := projectProductResult.accountStructureId
		propertyStructureId := projectProductResult.propertyStructureId

		for a := 0; a < 10; a++ {
			genAccount, err := dataGeneration.GenerateSingleAccount(groupIds)
			if err != nil {
				handleAppError(err, Cannot_Continue_Procedure)
			}

			joinedAccount := newJoinedStructureAccount(
				projectId,
				accountStructureId,
				propertyStructureId,
				groupIds,
				genAccount,
			)

			wq.addJob(newAccountWorkQueueJob(
				client,
				joinedAccount.projectId,
				joinedAccount.accountStructureId,
				joinedAccount.propertyStructureId,
				joinedAccount.groupIds,
				joinedAccount.account,
			))
		}
	}

	return publishingListeners
}

func concurrencyCoordinator(
	propertiesWorkQueue propertiesWorkQueue,
	accountWorkQueue *accountWorkQueue,
	progressBarNotifier chan bool,
	propertyWorkQueueDone chan bool,
	accountWorkQueueDone chan bool,
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
			case <-accountWorkQueue.jobDoneQueue:
				progressBarNotifier <- true
				reporter.AddAccount()
				workQueueTimeout = time.After(5 * time.Second)
			case <-workQueueTimeout:
				close(propertyWorkQueueDone)
				close(accountWorkQueueDone)
				return
			}
		}
	}()
}

func mergeDoneQueues(accountQueue chan bool, propertyQueue chan bool) chan bool {
	done := make(chan bool)

	go func() {
		<-accountQueue
		<-propertyQueue

		done <- true
	}()

	return done
}
