package main

import (
	"fmt"
	"sync"
)

/**
WARNING: THIS IS A DESTRUCTIVE COMMAND. IN CASE OF CERTAIN ERRORS, IT MIGHT DESTROY ALL THE DATA THAT YOU HAVE
		 IN THE DATABASE. USE WITH CAUTION!!!
*/

func main() {
	loadEnv()
	runDb()
	reactToFlags()

	anonymousClient := createAnonymousClient()
	authenticatedClient := preSeedAuthAndSetup(anonymousClient)

	printers["info"].Println("Creating projects")
	projects := generateProjects(authenticatedClient)

	printers["info"].Println("Creating project data with groups, Account(s) and Property(s)")
	fmt.Println("")

	progressBarNotifier := generateProgressBar(1000)

	wg := sync.WaitGroup{}
	wg.Add(len(projects))
	for _, p := range projects {
		projectId := p.id

		go func(projectId string) {
			defer wg.Done()

			handleHttpError(createGroups(authenticatedClient, projectId), nil)
			handleHttpError(createMapStructure(authenticatedClient, projectId, "Accounts"), nil)
			handleHttpError(createListStructure(authenticatedClient, projectId, "Properties"), nil)

			generatedAccounts, err := generateAccountStructureData("Accounts")
			if err != nil {
				handleAppError(err, Cannot_Continue_Procedure)
			}

			for _, genAccount := range generatedAccounts {
				handleHttpError(addToMap(authenticatedClient, projectId, genAccount.name, genAccount.variable, genAccount.references, genAccount.imagePaths), nil)
				progressBarNotifier <- true
			}
		}(projectId)
	}

	wg.Wait()
	close(progressBarNotifier)

	fmt.Println("")
	printers["success"].Println("Seed is successful!")
	fmt.Println("")
}
