package pagination

import (
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

var validOrderByFields = []string{
	"name",
	"created_at",
	"updated_at",
	"behaviour",
	"index",
}

var validFields = []string{
	"value",
	"metadata",
	"groups",
}

var initialReturnFields = []string{
	"name",
	"created_at",
	"updated_at",
	"behaviour",
	"index",
}

type Model struct {
	ProjectID        string
	StructureID      string
	StructureType    string
	ParentVariableID string
	Locales          []string

	Limit          int
	Page           int
	Search         string
	Filters        map[string]string
	OrderBy        string
	Behaviour      string
	OrderDirection string
	Groups         []string
	Fields         []string
}

func NewModel(
	projectId,
	structureType,
	parentVariableId string,
	locales []string,
	structureId,
	orderBy,
	search,
	direction string,
	limit, page int,
	groups []string,
	filters map[string]string,
	behaviour string,
	fields []string,
) Model {
	return Model{
		ProjectID:        projectId,
		Locales:          locales,
		Search:           search,
		StructureType:    structureType,
		ParentVariableID: parentVariableId,
		StructureID:      structureId,
		OrderBy:          orderBy,
		Page:             page,
		Filters:          filters,
		Behaviour:        behaviour,
		OrderDirection:   direction,
		Limit:            limit,
		Groups:           groups,
		Fields:           fields,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID":        a.ProjectID,
		"parentVariableId": a.ParentVariableID,
		"orderBy":          a.OrderBy,
		"page":             a.Page,
		"validFields":      a.Fields,
		"limit":            a.Limit,
		"behaviour":        a.Behaviour,
		"structureType":    a.StructureType,
		"direction":        a.OrderDirection,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("parentVariableId", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("structureType", validation.By(func(value interface{}) error {
				t := value.(string)

				if t != "map" && t != "list" {
					return errors.New("Structure type can be either 'map' or 'list'")
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
			validation.Key("validFields", validation.By(func(value interface{}) error {
				v := value.([]string)
				if len(v) == 0 {
					return nil
				}

				for _, r := range v {
					if !sdk.Includes(validFields, r) {
						return errors.New(fmt.Sprintf("Invalid return field. Valid return fields are '%s'. '%s' given", strings.Join(validFields, ","), r))
					}
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
