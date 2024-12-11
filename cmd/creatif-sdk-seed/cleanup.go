package main

import (
	"creatif-sdk-seed/errorHandler"
	"fmt"
)

func doOrderedCleanup() {
	printNewlineSandwich(printers["info"], "Starting cleanup...")
	errorHandler.CompleteCleanup()
	printers["success"].Println("Cleanup successful!")
	fmt.Println("")
}
