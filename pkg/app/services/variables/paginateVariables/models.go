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
	Locales   []string

	Limit          int
	Page           int
	Search         string
	Behaviour      string
	Filters        map[string]string
	OrderBy        string
	OrderDirection string
	Groups         []string
}

func NewModel(projectId string, locales []string, orderBy, search, direction string, limit, page int, groups []string, behaviour string, filters map[string]string) Model {
	return Model{
		ProjectID:      projectId,
		Locales:        locales,
		Search:         search,
		OrderBy:        orderBy,
		Behaviour:      behaviour,
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
		"locales":   a.Locales,
		"orderBy":   a.OrderBy,
		"behaviour": a.Behaviour,
		"page":      a.Page,
		"limit":     a.Limit,
		"direction": a.OrderDirection,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectId", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("locales", validation.By(func(value interface{}) error {
				t := value.([]string)

				if len(t) == 0 {
					return nil
				}

				for _, l := range t {
					if !locales.ExistsByAlpha(l) {
						return errors.New(fmt.Sprintf("Locale '%s' does not exist.", l))
					}
				}

				return nil
			})),
			validation.Key("behaviour", validation.By(func(value interface{}) error {
				v := value.(string)
				if v == "" {
					return nil
				}

				if v != "modifiable" && v != "readonly" {
					return errors.New("Behaviour can be only 'modifiable' and 'readonly'")
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
