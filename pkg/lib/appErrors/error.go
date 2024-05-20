package appErrors

import (
	"encoding/json"
	"errors"
	"fmt"
)

type AppError[T any] interface {
	AddError(method string, err error) AppError[T]
	Data() T
	Type() int
	StackTrace() string
	Error() string
	JSON() ([]byte, error)
}

type appError[T any] struct {
	data   T
	t      int
	errors []error
	stack  error
}

func (e *appError[T]) AddError(method string, err error) AppError[T] {
	if err != nil {
		e.errors = append(e.errors, fmt.Errorf("%s: %w", method, err))
	} else {
		e.errors = append(e.errors, fmt.Errorf("%s", method))
	}

	return e
}

func (e *appError[T]) Data() T {
	return e.data
}

func (e *appError[T]) Type() int {
	return e.t
}

func (e *appError[T]) StackTrace() string {
	str := ""
	counter := 1
	for i := len(e.errors) - 1; i >= 0; i-- {
		err := e.errors[i]

		str += fmt.Sprintf("%d. ", counter)
		if err != nil {
			str += fmt.Sprintf("%s\n", err.Error())
		}

		counter++
	}

	return str
}

func (e *appError[T]) JSON() ([]byte, error) {
	d, err := json.Marshal(e.Data())
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (e *appError[T]) Error() string {
	if len(e.errors) > 1 {
		return e.errors[0].Error()
	}

	return "No error provided"
}

func NewValidationError(data map[string]string) AppError[map[string]string] {
	errs := make([]error, 0)
	errs = append(errs, errors.New("Failed validation"))

	return &appError[map[string]string]{
		data:   data,
		t:      VALIDATION_ERROR,
		errors: errs,
		stack:  nil,
	}
}

func NewDatabaseError(err error) AppError[struct{}] {
	errs := make([]error, 0)
	errs = append(errs, err)

	return &appError[struct{}]{
		data:   struct{}{},
		t:      DATABASE_ERROR,
		errors: errs,
		stack:  nil,
	}
}

func NewUnexpectedError(err error) AppError[struct{}] {
	errs := make([]error, 0)
	errs = append(errs, err)

	return &appError[struct{}]{
		data:   struct{}{},
		t:      UNEXPECTED_ERROR,
		errors: errs,
		stack:  nil,
	}
}

func NewUserUnconfirmedError(err error) AppError[struct{}] {
	errs := make([]error, 0)
	errs = append(errs, err)

	return &appError[struct{}]{
		data:   struct{}{},
		t:      USER_UNCOFIRMED,
		errors: errs,
		stack:  nil,
	}
}

func NewNotFoundError(err error) AppError[struct{}] {
	errs := make([]error, 0)
	errs = append(errs, err)

	return &appError[struct{}]{
		data:   struct{}{},
		t:      NOT_FOUND_ERROR,
		errors: errs,
		stack:  nil,
	}
}

func NewApplicationError(err error) AppError[struct{}] {
	errs := make([]error, 0)
	errs = append(errs, err)

	return &appError[struct{}]{
		data:   struct{}{},
		t:      APPLICATION_ERROR,
		errors: errs,
		stack:  nil,
	}
}

func NewAuthorizationError(err error) AppError[struct{}] {
	errs := make([]error, 0)
	errs = append(errs, err)

	return &appError[struct{}]{
		data:   struct{}{},
		t:      AUTHORIZATION_ERROR,
		errors: errs,
		stack:  nil,
	}
}

func NewAuthenticationError(err error) AppError[struct{}] {
	errs := make([]error, 0)
	errs = append(errs, err)

	return &appError[struct{}]{
		data:   struct{}{},
		t:      AUTHENTICATION_ERROR,
		errors: errs,
		stack:  nil,
	}
}
