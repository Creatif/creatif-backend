package paginateVariables

import (
	"creatif/pkg/app/services/locales"
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
	"behaviour",
}

type Model struct {
	ProjectID string
	Locale    string

	Limit          int
	Page           int
	Search         string
	Filters        map[string]string
	OrderBy        string
	OrderDirection string
	Groups         []string
}

func NewModel(projectId, locale, orderBy, direction string, limit, page int, groups []string, filters map[string]string) Model {
	return Model{
		ProjectID:      projectId,
		Locale:         locale,
		OrderBy:        orderBy,
		Page:           page,
		Filters:        filters,
		OrderDirection: direction,
		Limit:          limit,
		Groups:         groups,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectId": a.ProjectID,
		"locale":    a.Locale,
		"orderBy":   a.OrderBy,
		"page":      a.Page,
		"limit":     a.Limit,
		"direction": a.OrderDirection,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectId", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("locale", validation.Required, validation.By(func(value interface{}) error {
				t := value.(string)

				if !locales.ExistsByAlpha(t) {
					return errors.New(fmt.Sprintf("Locale '%s' does not exist.", t))
				}

				return nil
			})),
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
