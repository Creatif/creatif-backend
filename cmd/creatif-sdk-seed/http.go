package main

import "net/http"

type result struct {
	response  *http.Response
	error     error
	status    int
	ok        bool
	procedure string
}

func (r result) Response() *http.Response {
	return r.response
}

func (r result) Error() error {
	return r.error
}

func (r result) Status() int {
	return r.status
}

func (r result) Procedure() string {
	return r.procedure
}

func (r result) Ok() bool {
	return r.ok
}

func newHttpResult(response *http.Response, error error, status int, ok bool, procedure string) httpResult {
	return result{
		response:  response,
		error:     error,
		status:    status,
		ok:        ok,
		procedure: procedure,
	}
}
