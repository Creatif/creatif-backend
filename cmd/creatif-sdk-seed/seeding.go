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

type createdClientVariable struct {
	projectId           string
	clientId            string
	propertyStructureId string
	groupIds            []string
}

type projectProduct struct {
	projectId           string
	groupIds            []string
	clientStructureId   string
	propertyStructureId string
	managerStructureId  string
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

func newProjectProduct(projectId, clientStructureId, propertyStructureId, managerStructureId string, groupIds []string) projectProduct {
	return projectProduct{
		projectId:           projectId,
		groupIds:            groupIds,
		clientStructureId:   clientStructureId,
		managerStructureId:  managerStructureId,
		propertyStructureId: propertyStructureId,
	}
}

func projectProducer(client *http.Client, numOfProjects int) []projectProduct {
	producers := make([]projectProduct, numOfProjects)

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
		managerStructureId := createManagersStructureAndReturnId(client, projectId)
		propertyStructureId := createPropertiesStructureAndReturnID(client, projectId)

		producers[i] = newProjectProduct(projectId, clientStructureId, propertyStructureId, managerStructureId, groupIds)
	}

	return producers
}

func createClientsAndManagers(
	clientNum int,
	client *http.Client,
	projects []projectProduct,
	reporter *reporter,
) []createdClientVariable {
	clientVariables := make([]createdClientVariable, len(projects)*clientNum)
	clientVariablesCount := 0
	printers["info"].Println("Creating clients and managers")
	for _, projectProductResult := range projects {
		reporter.AddProjectID(projectProductResult.projectId)

		groupIds := projectProductResult.groupIds
		projectId := projectProductResult.projectId
		clientStructureId := projectProductResult.clientStructureId
		managerStructureId := projectProductResult.managerStructureId
		propertyStructureId := projectProductResult.propertyStructureId

		projectClientVariables := make([]createdClientVariable, clientNum)
		for a := 0; a < clientNum; a++ {
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

			clientId := addToMapAndGetClientId(
				client,
				projectId,
				joinedClient.clientStructureId,
				joinedClient.client,
			)

			createdVariable := createdClientVariable{
				projectId:           projectId,
				clientId:            clientId,
				propertyStructureId: propertyStructureId,
				groupIds:            groupIds,
			}

			clientVariables[clientVariablesCount] = createdVariable
			projectClientVariables[a] = createdVariable

			reporter.AddClient()
			clientVariablesCount++
		}

		clientConnections := make([]dataGeneration.ClientConnection, len(projectClientVariables))
		for i, c := range projectClientVariables {
			clientConnections[i] = dataGeneration.ClientConnection{
				StructureType: "map",
				VariableID:    c.clientId,
			}
		}

		for i := 0; i < 100; i++ {
			manager, err := dataGeneration.GenerateSingleManager(groupIds, clientConnections)
			if err != nil {
				errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
			}

			errorHandler.HandleHttpError(createManager(
				client,
				managerStructureId,
				projectId,
				manager.Variable,
				manager.Connections,
				manager.ImagePaths,
			))
		}
	}

	return clientVariables
}

func startSeeding(
	client *http.Client,
	wq *propertiesWorkQueue,
	propertiesPerStatus int,
	clientVariables []createdClientVariable,
) {
	printers["info"].Println("Creating properites\n")
	for _, clientVariable := range clientVariables {
		propertiesGen := dataGeneration.NewPropertiesGenerator()

		for {
			newSequence, ok := propertiesGen.Generate()

			if !ok {
				break
			}

			for a := 0; a < propertiesPerStatus; a++ {
				singleProperty, err := dataGeneration.GenerateSingleProperty(
					clientVariable.clientId,
					newSequence.Locale,
					newSequence.PropertyStatus,
					newSequence.PropertyType,
					clientVariable.groupIds,
				)
				if err != nil {
					errorHandler.HandleAppError(err, Cannot_Continue_Procedure)
				}

				wq.addJob(newPropertyWorkQueueJoby(
					client,
					clientVariable.projectId,
					clientVariable.propertyStructureId,
					singleProperty.Variable,
					singleProperty.Connections,
					singleProperty.ImagePaths,
				))
			}
		}
	}
}

func concurrencyCoordinator(
	propertiesWorkQueue propertiesWorkQueue,
	propertyWorkQueueDone chan bool,
	reporter *reporter,
) {
	go func() {
		workQueueTimeout := time.After(5 * time.Second)
		doneProperties := 0
		for {
			select {
			case <-propertiesWorkQueue.jobDoneQueue:
				reporter.AddProperty()
				doneProperties++

				if doneProperties%1000 == 0 {
					fmt.Printf("Processed %d properties\n", doneProperties)
				}

				workQueueTimeout = time.After(5 * time.Second)
			case <-workQueueTimeout:
				propertyWorkQueueDone <- true
				close(propertyWorkQueueDone)
				return
			}
		}
	}()
}
