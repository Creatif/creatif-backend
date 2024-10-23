package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"sync"
)

/**
WARNING: THIS IS A DESTRUCTIVE COMMAND. IN CASE OF CERTAIN ERRORS, IT MIGHT DESTROY ALL THE DATA THAT YOU HAVE
		 IN THE DATABASE. USE WITH CAUTION!!!
*/

func main() {
	loadEnv()
	runDb()

	shouldJustCleanup := len(os.Args) == 2 && os.Args[1] == "--cleanup"
	if shouldJustCleanup {
		doOrderedCleanup()
		os.Exit(0)
	}

	successColor := color.New(color.FgGreen).Add(color.Bold)

	anonymousClient := createAnonymousClient()
	email := "email@gmail.com"
	password := "password"

	if adminExists(anonymousClient).Ok() {
		printNewlineSandwich(successColor, "Admin already exists which means that the seed is there.\nIf it is not and this is a mistake, just delete the docker volume and try again.\nThis is fine and OK since this is a seed program to test the SDK.\nFeel free to abuse it.")
		return
	}

	printers["info"].Println("Creating admin and logging in")
	handleHttpError(createAdmin(anonymousClient, email, password), nil)

	authToken := extractAuthenticationCookie(handleHttpError(login(anonymousClient, email, password), nil))

	authenticatedClient := createAuthenticatedClient(authToken)

	printers["info"].Println("Creating projects")
	projects := generateProjects(authenticatedClient)

	printers["info"].Println("Creating project data with groups, Account(s) and Property(s)")

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
			}
		}(projectId)
	}

	wg.Wait()

	fmt.Println("")
	printers["success"].Println("Seed is successful!")
}
