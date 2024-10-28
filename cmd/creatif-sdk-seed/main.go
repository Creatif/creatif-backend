/**
WARNING: THIS IS A DESTRUCTIVE COMMAND. IN CASE OF CERTAIN ERRORS, IT MIGHT DESTROY ALL THE DATA THAT YOU HAVE
		 IN THE DATABASE. USE WITH CAUTION!!!

IMPORTANT:
This seed actually uploads images. Every account gets one image and every property gets 3 images. It would be wise
to from time to time, just delete the 'var' and 'public' directories because they might get very large if you execute
this function over and over again.

This program cannot start if you don't have the server up, so make sure that you open up a new terminal tab, hit 'docker compose up' on the main project
and only then execute this command.

This command seeds the initial application with seed data from real estate project. It has two structures: Accounts and
Properties. Accounts is a map and Properties is a list. It generates five projects with those structure. Each project has
200 Account maps and 1000 (one thousand) Properties in 5 different locales. That means that this command will generate 5200
"entities" per project. There will be 5 projects so 26 thousand "entities" will be created in total.

This command will be used to test public SDKs. For now, there is only javascript SDK but hopefully, there will be more.

If you try to execute this command more than once, it will not give that to you. Since the application can have a single admin (for now),
you cannot create another admin, therefor the program will tell you that the app is already seeded.

There is nothing special about this program. Just cd into this directory and run 'go run .' and that is it. Calling this program without
any flags will create five projects by default.

Flags:
--cleanup
    This flag will completely destroy all data in the database. USE WITH CAUTION!!! If you use
    this flag is the only thing that will be done even if you used other flags i.e. it will ignore all other flags.
--regenerate
    This will do what --cleanup does but will run other commands. Basically, you tell the program to wipe
    the database out start over
--projects={\d} (default is 5)
    For how many projects should it seed the application. More will be slower.
--help
    If used, will output help. If used with any other flags, it will ignore them and just print help, i.e. it will
    ignore all other flags.

Credentials:

Email: mariofake@gmail.com
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

	report := newReporter()

	anonymousClient := createAnonymousClient()
	authenticatedClient := preSeedAuthAndSetup(anonymousClient)

	propertiesWorkQueue := newListWorkQueue(50, 10)
	propertyWorkQueueDone := propertiesWorkQueue.start()

	accountWorkQueue := newMapWorkQueue(50, 10, propertiesWorkQueue)
	accountWorkQueueDone := accountWorkQueue.start()

	fmt.Printf("Seeding...\n")
	fmt.Println("")

	numOfAllOperations := (numOfProjects * 10) * 600
	progressBarNotifier, progressBarDone := generateProgressBar(numOfAllOperations)

	projectProducerListeners := projectProducer(authenticatedClient, numOfProjects)
	projectPublishingListeners := accountProducer(authenticatedClient, projectProducerListeners, accountWorkQueue, report)

	concurrencyCoordinator(
		propertiesWorkQueue,
		accountWorkQueue,
		progressBarNotifier,
		propertyWorkQueueDone,
		accountWorkQueueDone,
		report,
	)

	<-mergeDoneQueues(accountWorkQueueDone, propertyWorkQueueDone)
	progressBarDone <- true

	fmt.Println("")
	printers["info"].Println("Publishing projects...")
	publishProjects(authenticatedClient, projectPublishingListeners)
	printers["info"].Println("Projects published")

	fmt.Println("")
	printers["success"].Println("Seed is successful!")
	report.Report()
	fmt.Println("")
}