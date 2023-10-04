package paginateVariables

import (
	"creatif/pkg/app/services/paginateVariables/pagination"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

type Model struct {
	PaginationID string
	// field for ORDER_BY clause
	Field string
	// DESC or ASC
	OrderBy string
	// forward or backwards
	Direction string
	Groups    []string
	Limit     int
	ProjectID string
}

func NewModel(projectId, paginationId, field, orderBy, direction string, limit int, groups []string) Model {
	return Model{
		ProjectID:    projectId,
		PaginationID: paginationId,
		Field:        field,
		OrderBy:      orderBy,
		Direction:    direction,
		Limit:        limit,
		Groups:       groups,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"paginationId": a.PaginationID,
		"field":        a.Field,
		"orderBy":      a.OrderBy,
		"direction":    a.Direction,
		"limit":        a.Limit,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("paginationId", validation.When(a.PaginationID != "", validation.Required, validation.RuneLength(26, 26))),
			validation.Key("field", validation.Required, validation.RuneLength(1, 50)),
			validation.Key("orderBy", validation.Required, validation.By(func(value interface{}) error {
				t := value.(string)
				orderBy := strings.ToUpper(t)

				if orderBy != pagination.DESC && orderBy != pagination.ASC {
					return errors.New(fmt.Sprintf("orderBy field can be either '%s' or '%s'", pagination.DESC, pagination.ASC))
				}

				a.OrderBy = orderBy

				return nil
			})),
			validation.Key("direction", validation.Required, validation.By(func(value interface{}) error {
				t := value.(string)

				if t != pagination.DIRECTION_FORWARD && t != pagination.DIRECTION_BACKWARDS {
					return errors.New(fmt.Sprintf("direction field can be either '%s' or '%s'", pagination.DIRECTION_FORWARD, pagination.DIRECTION_BACKWARDS))
				}

				return nil
			})),
			validation.Key("limit", validation.Required, validation.By(func(value interface{}) error {
				t := value.(int)

				if t <= 0 {
					return errors.New("Limit must be greater than 0 (zero)")
				}

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
