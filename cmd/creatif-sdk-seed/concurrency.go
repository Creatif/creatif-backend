package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type joinedStructureAccount struct {
	projectId           string
	groupIds            []string
	accountId           string
	accountStructureId  string
	propertyStructureId string

	account account
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

	account account,
) joinedStructureAccount {
	return joinedStructureAccount{
		projectId:           projectId,
		groupIds:            groupIds,
		accountStructureId:  accountStructureId,
		propertyStructureId: propertyStructureId,
		account:             account,
	}
}

func newAccountProduct(projectId, accountId, propertyStructureId string, groupIds []string) accountProduct {
	return accountProduct{
		projectId:           projectId,
		accountId:           accountId,
		groupIds:            groupIds,
		propertyStructureId: propertyStructureId,
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
			handleHttpError(createProject(client, projectName), func(res *http.Response) error {
				var m map[string]interface{}
				b, err := io.ReadAll(res.Body)
				defer res.Body.Close()
				if err != nil {
					return err
				}

				if err := json.Unmarshal(b, &m); err != nil {
					return err
				}

				if res.StatusCode < 200 && res.StatusCode > 299 {
					return errors.New(fmt.Sprintf("Creating project failed with status %d and body %s", res.StatusCode, string(b)))
				}

				projectId = m["id"].(string)

				groupIds := createGroupsAndGetGroupIds(client, projectId)
				accountStructureId := createAccountStructureAndReturnID(client, projectId)
				propertyStructureId := createPropertiesStructureAndReturnID(client, projectId)

				producers[i] <- newProjectProduct(projectId, accountStructureId, propertyStructureId, groupIds)

				return nil
			})
		}
	}()

	return producers
}

func accountProducer(client *http.Client, projectProducers []chan projectProduct, wq mapWorkQueue) {
	allAccounts := make([]joinedStructureAccount, 0)
	for _, projectProducer := range projectProducers {
		projectProductResult := <-projectProducer
		groupIds := projectProductResult.groupIds
		projectId := projectProductResult.projectId
		accountStructureId := projectProductResult.accountStructureId
		propertyStructureId := projectProductResult.propertyStructureId
		generatedAccounts, err := generateAccountStructureData(groupIds)
		if err != nil {
			handleAppError(err, Cannot_Continue_Procedure)
		}

		joinedAccounts := make([]joinedStructureAccount, 0)
		for _, genAccount := range generatedAccounts {
			joinedAccounts = append(joinedAccounts, newJoinedStructureAccount(
				projectId,
				accountStructureId,
				propertyStructureId,
				groupIds,
				genAccount,
			))
		}

		allAccounts = append(allAccounts, joinedAccounts...)
	}

	for _, joinedAccount := range allAccounts {
		wq.addJob(newMapWorkQueueJob(
			client,
			joinedAccount.projectId,
			joinedAccount.accountStructureId,
			joinedAccount.propertyStructureId,
			joinedAccount.groupIds,
			joinedAccount.account,
		))
	}
}

func concurrencyCoordinator(
	propertiesWorkQueue listWorkQueue,
	accountWorkQueue mapWorkQueue,
	progressBarNotifier chan bool,
	numOfAllOperations int,
	propertyWorkQueueDone chan bool,
	accountWorkQueueDone chan bool,
) {
	go func() {
		operations := 0
		for {
			select {
			case <-propertiesWorkQueue.jobDoneQueue:
				progressBarNotifier <- true
				operations++
			case <-accountWorkQueue.jobDoneQueue:
				progressBarNotifier <- true
				operations++
			default:
				if numOfAllOperations == operations {
					close(propertyWorkQueueDone)
					close(accountWorkQueueDone)
					return
				}
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
