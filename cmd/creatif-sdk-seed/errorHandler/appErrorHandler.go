package errorHandler

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

const URL = "http://localhost:3002/api/v1"

const Cannot_Continue_Procedure = "cannot_continue"
const User_Error_Procedure = "user_error_can_continue"

var printers = map[string]*color.Color{
	"error": color.New(color.FgRed).Add(color.Bold),
}

const (
	PROJECT_TABLE          = "projects"
	USERS_TABLE            = "users"
	MAP_VARIABLES          = "map_variables"
	VARIABLE_MAP           = "maps"
	LIST_TABLE             = "lists"
	GROUPS_TABLE           = "groups"
	VARIABLE_GROUPS_TABLE  = "variable_groups"
	CONNECTIONS_TABLE      = "connections"
	LIST_VARIABLES_TABLE   = "list_variables"
	PUBLISHED_LISTS_TABLE  = "published_lists"
	PUBLISHED_MAPS_TABLE   = "published_maps"
	VERSION_TABLE          = "versions"
	FILE_TABLE             = "files"
	PUBLISHED_FILES_TABLE  = "published_files"
	PUBLISHED_GROUPS_TABLE = "published_groups"
)

func printNewlineSandwich(printer *color.Color, print string) {
	printer.Println(print)
}

func HandleAppError(err error, flag string) {
	if flag == Cannot_Continue_Procedure {
		printNewlineSandwich(printers["error"], fmt.Sprintf("An app error occurred: %s. The program is forced to quit.", err.Error()))
		printNewlineSandwich(printers["error"], "The program is forced to clean up everything that happened up until now.\nThat means a complete database wipe out.\nRun this command again for a clean start.")
		printNewlineSandwich(printers["error"], "IMPORTANT: cleanup truncates every table it the database but does not check if it errors.\nIt is perfectly fine to delete the docker volume to start again.")
		fmt.Println("")
		CompleteCleanup()
		os.Exit(1)
	}

	if err.Error() == "invalid_num_images" {
		printNewlineSandwich(printers["error"], fmt.Sprintf("An app error occurred: %s. The program is forced to quit.", err.Error()))
		printNewlineSandwich(printers["error"], "The program is forced to clean up everything that happened up until now.\nThat means a complete database wipe out.\nRun this command again for a clean start.")
		printNewlineSandwich(printers["error"], "IMPORTANT: cleanup truncates every table it the database but does not check if it errors.\nIt is perfectly fine to delete the docker volume to start again.")
		fmt.Println("")
		CompleteCleanup()
		os.Exit(1)
	}
	// ignore the error, it is not serious enough
}
