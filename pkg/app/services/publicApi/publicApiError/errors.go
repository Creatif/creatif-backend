package publicApiError

import (
	"encoding/json"
	"net/http"
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
	s := http.StatusInternalServerError
	if status == ValidationError {
		s = http.StatusUnprocessableEntity
	} else if status == NotFoundError {
		s = http.StatusNotFound
	} else if status == ApplicationError {
		s = http.StatusBadRequest
	}

	return appError{
		data: map[string]interface{}{
			"call":     call,
			"messages": messages,
			"status":   s,
		},
		status: status,
	}
}
