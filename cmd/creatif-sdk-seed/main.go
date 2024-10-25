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
	"os"
)

func main() {
	loadEnv()
	runDb()
	numOfProjects, err := processFlags()
	if err != nil {
		printNewlineSandwich(printers["error"], err.Error())
		os.Exit(1)
	}

	anonymousClient := createAnonymousClient()
	authenticatedClient := preSeedAuthAndSetup(anonymousClient)

	propertiesWorkQueue := newListWorkQueue(50, 10)
	propertyWorkQueueDone := propertiesWorkQueue.start()

	accountWorkQueue := newMapWorkQueue(50, 10, propertiesWorkQueue)
	accountWorkQueueDone := accountWorkQueue.start()

	fmt.Println("Seeding...")
	fmt.Println("")

	numOfAllOperations := (numOfProjects * 10) * 300
	progressBarNotifier, progressBarDone := generateProgressBar(numOfAllOperations)

	projectProducerListeners := projectProducer(authenticatedClient, numOfProjects)
	accountProducer(authenticatedClient, projectProducerListeners, accountWorkQueue)

	concurrencyCoordinator(
		propertiesWorkQueue,
		accountWorkQueue,
		progressBarNotifier,
		numOfAllOperations,
		propertyWorkQueueDone,
		accountWorkQueueDone,
	)

	<-mergeDoneQueues(accountWorkQueueDone, propertyWorkQueueDone)
	progressBarDone <- true

	fmt.Println("")
	printers["success"].Println("Seed is successful!")
	fmt.Println("")
}
