package http

import (
	"net/http"
)

type Request struct {
	Headers map[string]string
	Url     string
	Method  string
	Body    []byte
}

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

func NewHttpResult(response *http.Response, error error, status int, ok bool, procedure string) HttpResult {
	return result{
		response:  response,
		error:     error,
		status:    status,
		ok:        ok,
		procedure: procedure,
	}
}
