package main

import (
	"fmt"
)

func completeCleanup() {
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", PROJECT_TABLE))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", MAP_VARIABLES))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", VARIABLE_MAP))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", LIST_TABLE))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", LIST_VARIABLES_TABLE))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE app.%s CASCADE", USERS_TABLE))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", REFERENCE_TABLES))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", GROUPS_TABLE))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", VARIABLE_GROUPS_TABLE))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", PUBLISHED_LISTS_TABLE))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", PUBLISHED_MAPS_TABLE))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", VERSION_TABLE))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", PUBLISHED_REFERENCES_TABLE))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE declarations.%s CASCADE", FILE_TABLE))
	Gorm().Exec(fmt.Sprintf("TRUNCATE TABLE published.%s CASCADE", PUBLISHED_FILES_TABLE))
}

func doOrderedCleanup() {
	printNewlineSandwich(printers["info"], "Starting cleanup...")
	completeCleanup()
	printers["success"].Println("Cleanup successful!")
	fmt.Println("")
}
