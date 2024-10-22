package main

import "net/http"

type httpResult interface {
	Response() *http.Response
	Procedure() string
	Error() error
	Status() int
	Ok() bool
}
