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

	/**
		These are the work queues where Accounts and Properties are created. Each of them has 50 workers that produce
		50 workers and each of them produce 50 channels that others can send work to. Very important thing to note is that
		these are not generic work queues. Both of them are specialized to produce either Account or a Property.

		Below, accountProducer sends a job to mapWorkQueue for a Account to be created. After that, it creates a job for
		the listWorkQueue for a property to be created. That way, they do not block each other. Every worker has its own
		channel that is buffered. That is the second argument.

		These queues can be started before any job is sent to them.

	    start() function returns a to be closed when all work is done. These workers will ALWAYS have work to be done. It
		is important to know that. For more information how these are used, take a look below to the comment above
		concurrencyCoordinator() function.
	*/
	propertiesWorkQueue := newListWorkQueue(60, 60)
	propertyWorkQueueDone := propertiesWorkQueue.start()
	accountWorkQueue := newMapWorkQueue(60, 60, propertiesWorkQueue)
	accountWorkQueueDone := accountWorkQueue.start()

	fmt.Printf("Seeding...\n")
	fmt.Println("")

	numOfAllOperations := (numOfProjects * 10) * 600
	progressBarNotifier, progressBarDone := generateProgressBar(numOfAllOperations)

	/**
	projectProducer produces as many producer channels as numOfProjects. You can listen to these
	producer channels when a project has been created. This is of course done in 1 goroutine per
	project. For example, if numOfProjects is 5, there will be 5 goroutines that will create projects
	and this function will return 5 producer channels to listen to when a project has been created.
	*/
	projectProducerListeners := projectProducer(authenticatedClient, numOfProjects)
	/**
	accountProducer listens to when a project is created and sends it to mapWorkQueue. More on mapWorkQueue, just scroll
	up. Not complicated, for every project, there are 10 Accounts created and send to mapWorkQueue which is just another
	work queue.
	*/
	projectPublishingListeners := accountProducer(authenticatedClient, projectProducerListeners, accountWorkQueue, report)

	/**
	    This function blocks until all work queues are done with their jobs. But this is not the join point. Every work queue exposes
		a channel that signals when a job is done. A reporter is here just to write to stdout to the user of this program
		a user-friendly message how many jobs have been done.

	    propertyWorkQueueDone and accountWorkQueueDone are channels that signal to these work queues that their job is done,
		and they can be garbage collected. The way they know that is a timeout. If any of the worker has nothing to do
		for more than 2 seconds, the queues are closed and garbage collected. This is OK because the program is made so that,
		if there is work to be done, there will always be work. If there is not work in both queues, that means that we have
		seeded everything that needs to be seeded and work queues can be closed.
	*/
	concurrencyCoordinator(
		propertiesWorkQueue,
		accountWorkQueue,
		progressBarNotifier,
		propertyWorkQueueDone,
		accountWorkQueueDone,
		report,
	)

	/**
	This is the join point to block the work queues until they are finished. close() function on these two channels
	is called in concurrencyCoordinator() after all the work is done.
	*/
	<-mergeDoneQueues(accountWorkQueueDone, propertyWorkQueueDone)
	progressBarDone <- true

	fmt.Println("")
	printers["info"].Println("Publishing projects...")
	publishProjects(authenticatedClient, projectPublishingListeners)
	printers["info"].Println("Projects published")

	fmt.Println("")
	printers["success"].Println("Seed is successful!")
	report.Report()
	fmt.Println(fmt.Sprintf("Email: %s", Email))
	fmt.Println(fmt.Sprintf("Password: %s", Password))
	fmt.Println("")
}
