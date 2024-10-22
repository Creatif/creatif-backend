package main

import "net/http"

type httpResult interface {
	Response() *http.Response
	// this tells the program if the error, success or any other situation requires special attention.
	// an http call might be 200, but should not be, therefor a special procedure must be created to handle it.
	Procedure() string
	Error() error
	Status() int
	// Status means that a function that does an http call is successfull from its point of view, not from
	// success or failure of the http request. For example, adminExists() function will have ok == false if
	// the admin exists. Every function can interpret what to make "ok" for itself.
	Ok() bool
}
