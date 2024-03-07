package publicApiError

import (
	"encoding/json"
)

const NotFoundError = 1
const DatabaseError = 2
const ApplicationError = 3
const ValidationError = 4

type ErrorData map[string]interface{}

type PublicApiError interface {
	Data() ErrorData
	Error() string
	Status() int
}

type appError struct {
	data   ErrorData
	status int
}

func (e appError) Data() ErrorData {
	return e.data
}

func (e appError) Error() string {
	j, _ := json.Marshal(e.data)
	return string(j)
}

func (e appError) Status() int {
	return e.status
}

func NewError(call string, messages map[string]string, status int) PublicApiError {
	return appError{
		data: map[string]interface{}{
			"call":     call,
			"messages": messages,
			"status":   status,
		},
		status: status,
	}
}
