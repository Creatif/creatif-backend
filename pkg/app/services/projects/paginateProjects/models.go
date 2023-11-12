package paginateProjects

import (
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

var validOrderByFields []string = []string{
	"name",
	"created_at",
	"updated_at",
}

type Model struct {
	Limit          int
	Page           int
	Search         string
	OrderBy        string
	OrderDirection string
}

func NewModel(orderBy, search, direction string, limit, page int) Model {
	return Model{
		Search:         search,
		OrderBy:        orderBy,
		Page:           page,
		OrderDirection: direction,
		Limit:          limit,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"orderBy":   a.OrderBy,
		"page":      a.Page,
		"limit":     a.Limit,
		"direction": a.OrderDirection,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("orderBy", validation.By(func(value interface{}) error {
				t := value.(string)

				if !sdk.Includes(validOrderByFields, t) {
					return errors.New(fmt.Sprintf("Invalid order_by field. Valid order_by fields are: '%s'", strings.Join(validOrderByFields, ",")))
				}

				return nil
			})),
			validation.Key("page", validation.By(func(value interface{}) error {
				if a.Page < 1 {
					return errors.New("Page must be either the number 1 or greater than 1.")
				}

				return nil
			})),
			validation.Key("limit", validation.By(func(value interface{}) error {
				if a.Limit < 1 {
					return errors.New("Limit must be either the number 1 or greater than 1. Maximum value is 1000.")
				}

				if a.Limit > 1000 {
					return errors.New("Limit must be either the number 1 or greater than 1. Maximum value is 1000.")
				}

				return nil
			})),
			validation.Key("direction", validation.By(func(value interface{}) error {
				d := strings.ToUpper(a.OrderDirection)
				if d != "ASC" && d != "DESC" {
					return errors.New("Order direction must be either ASC or DESC.")
				}

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
