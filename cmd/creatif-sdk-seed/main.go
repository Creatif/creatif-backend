/**
WARNING: THIS IS A DESTRUCTIVE COMMAND. IN CASE OF CERTAIN ERRORS, IT MIGHT DESTROY ALL THE DATA THAT YOU HAVE
		 IN THE DATABASE. USE WITH CAUTION!!! IT SHOULD NOT BE USED TO SEED THE APP FOR USE IN THE creatif-ui-sdk.

This command seeds the initial application with seed data from real estate project. It has two structures: Accounts and
Properties. Accounts is a map and Properties is a list. It generates five projects with those structure. Each project has
200 Account maps and 1000 (one thousand) Properties in 5 different locales. That means that this command will generate 5200
"entities" per project. There will be 5 projects so 26 thousand "entities" will be created in total.

This command will be used to test public SDKs. For now, there is only javascript SDK but hopefully, there will be more.

If you try to execute this command more than once, it will not give that to you. Since the application can have a single admin (for now),
you cannot create another admin, therefor the program will tell you that the app is already seeded.

There is nothing special about this program. Just cd into this directory and run 'go run .' and that is it.

Flags:
--cleanup
    This flag will completely destroy all data in the database. USE WITH CAUTION!!!



Email: skrlecmario88@gmail.com
Password: password

I know that password is weak, but there is a plan to put a password strength into effect in the application
but until then, this is just fine.
*/

package main

import (
	"fmt"
	"sync"
)

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
		projectId := p.ID

		go func(projectId string) {
			defer wg.Done()

			groupIds := createGroupsAndGetGroupIds(authenticatedClient, projectId)
			accountId := createAccountStructureAndReturnID(authenticatedClient, projectId)
			createPropertiesStructureAndReturnID(authenticatedClient, projectId)

			generatedAccounts, err := generateAccountStructureData(groupIds)
			if err != nil {
				handleAppError(err, Cannot_Continue_Procedure)
			}

			for _, genAccount := range generatedAccounts {
				handleHttpError(addToMap(authenticatedClient, projectId, accountId, genAccount.variable, genAccount.references, genAccount.imagePaths), nil)
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
