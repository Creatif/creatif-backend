package main

import "log"

const URL = "http://localhost:3002/api/v1"

const Cannot_Continue_Procedure = "cannot_continue"
const User_Error_Procedure = "user_error_can_continue"

func handleError(result httpResult) {
	if result.Ok() {
		return
	}

	if !result.Ok() && result.Procedure() == User_Error_Procedure {
		return
	}

	err := result.Error()
	if err != nil {
		// TODO: Possible cleanup here
		log.Fatalln(err)
	}
}
