package main

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"os"
	"time"
)

const URL = "http://localhost:3002/api/v1"

const Cannot_Continue_Procedure = "cannot_continue"
const User_Error_Procedure = "user_error_can_continue"
const Can_Continue = "Everything is OK"

const Email = "skrlecmario88@gmail.com"
const Password = "password"

const (
	PROJECT_TABLE              = "projects"
	USERS_TABLE                = "users"
	MAP_VARIABLES              = "map_variables"
	VARIABLE_MAP               = "maps"
	LIST_TABLE                 = "lists"
	GROUPS_TABLE               = "groups"
	VARIABLE_GROUPS_TABLE      = "variable_groups"
	REFERENCE_TABLES           = "references"
	LIST_VARIABLES_TABLE       = "list_variables"
	PUBLISHED_LISTS_TABLE      = "published_lists"
	PUBLISHED_MAPS_TABLE       = "published_maps"
	VERSION_TABLE              = "versions"
	FILE_TABLE                 = "files"
	PUBLISHED_FILES_TABLE      = "published_files"
	PUBLISHED_REFERENCES_TABLE = "published_references"
)

var printers = map[string]*color.Color{
	"success": color.New(color.FgGreen).Add(color.Bold),
	"error":   color.New(color.FgRed).Add(color.Bold),
	"info":    color.New(color.FgWhite).Add(color.Bold),
}

/*
*
Handles the httpResult produced by function that perform http calls. How each part works,
take a look at the actual body of the function. There are plenty of comments.

responseFn is optional and can be used to extract values from a response body. If this function
returns an error, it will consider that as a terminal error, cleanup the database and exit. If the
error that this callback produced is not that important, do not return an error.

responseFn callback will be called is the http function has response.Ok == true or if there is an
error and the response is available.

In any case, if you return an error from this function, CLEANUP OF THE SYSTEM WILL HAPPEN.
*/
func handleHttpError(result httpResult, responseFn func(res *http.Response) error) httpResult {
	// result is ok (>=200 <= 299). If the result is OK, then nothing else
	// matters and there is no error or no need to check the procedure
	if result.Ok() {
		if responseFn != nil {
			if err := responseFn(result.Response()); err != nil {
				printNewlineSandwich(printers["error"], fmt.Sprintf("Your callback produced an error: %s. The program is forced to quit.", err.Error()))
				printNewlineSandwich(printers["error"], "The program is forced to clean up everything that happened up until now.\nThat means a complete database wipe out.\nRun this command again for a clean start.")
				printNewlineSandwich(printers["error"], "IMPORTANT: cleanup truncates every table it the database but does not check if it errors.\nIt is perfectly fine to delete the docker volume to start again.")
				completeCleanup()
				os.Exit(1)
			}
		}

		return result
	}

	// this is also ok, since some requests can fail for various reasons like DOS the server.
	// This SHOULD NOT happen and the program is made for this not to happen, but if it does, it's fine.
	// In this case, the program will sleep for 2 seconds for the database to recover and continue again.
	if !result.Ok() && result.Procedure() == User_Error_Procedure {
		time.Sleep(2 * time.Second)
		return result
	}

	/**
	This is an error in the program, not the request/response lifecycle. For example, marshaling a map
	failed. This SHOULD NEVER happen and when it happens, it means there is a bug in the system.

	This program will delete everything that it has created so far and exit. The reader of this comment should
	investigate what the error was, where it originated and fix it since its most definitely a bug.
	*/
	err := result.Error()
	if err != nil && result.Procedure() == Cannot_Continue_Procedure {
		printNewlineSandwich(printers["error"], fmt.Sprintf("Something wrong happened here: %s. The program is forced to quit.", err.Error()))
		printNewlineSandwich(printers["error"], "The program is forced to clean up everything that happened up until now.\nThat means a complete database wipe out.\nRun this command again for a clean start.")
		printNewlineSandwich(printers["error"], "IMPORTANT: cleanup truncates every table it the database but does not check if it errors.\nIt is perfectly fine to delete the docker volume to start again.")

		res := result.Response()
		if res != nil {
			/**
			If there is a responseFn callback, that callback will be called for you to determine
			what to do with the response. If it is not provided, the response will be dumped to stdout
			for debug purposes.
			*/
			if responseFn != nil {
				if err := responseFn(result.Response()); err != nil {
					printNewlineSandwich(printers["error"], fmt.Sprintf("Your callback produced an error: %s. The program is forced to quit.", err.Error()))
					printNewlineSandwich(printers["error"], "The program is forced to clean up everything that happened up until now.\nThat means a complete database wipe out.\nRun this command again for a clean start.")
					printNewlineSandwich(printers["error"], "IMPORTANT: cleanup truncates every table it the database but does not check if it errors.\nIt is perfectly fine to delete the docker volume to start again.")
					completeCleanup()
					os.Exit(1)
				}
			} else {
				printNewlineSandwich(printers["info"], "There seems to be a response in this error. This program will print the first 128 characters of it just for the sake of debugging")
				bd, _ := io.ReadAll(result.Response().Body)
				defer result.Response().Body.Close()
				strBd := string(bd)
				if strBd == "" {
					printers["info"].Println("The response is empty")
				} else {
					fmt.Println(strBd)
				}
			}
		}

		completeCleanup()
		os.Exit(1)
	}

	return result
}

func handleAppError(err error, flag string) {
	if flag == Cannot_Continue_Procedure {
		printNewlineSandwich(printers["error"], fmt.Sprintf("An app error occurred: %s. The program is forced to quit.", err.Error()))
		printNewlineSandwich(printers["error"], "The program is forced to clean up everything that happened up until now.\nThat means a complete database wipe out.\nRun this command again for a clean start.")
		printNewlineSandwich(printers["error"], "IMPORTANT: cleanup truncates every table it the database but does not check if it errors.\nIt is perfectly fine to delete the docker volume to start again.")
		completeCleanup()
		os.Exit(1)
	}

	// ignore the error, it is not serious enough
}
