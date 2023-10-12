package getVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

var validFields = []string{
	"behaviour",
	"metadata",
	"groups",
	"value",
	"created_at",
	"updated_at",
}

type Model struct {
	// this can be project name
	Name        string `json:"name"`
	Fields      []string
	ProjectID   string `json:"projectID"`
	LocaleAlpha string `json:"localeAlpha"`

	validFields []string
}

func NewModel(projectId, name, localeAlpha string, fields []string) Model {
	if len(fields) == 0 {
		fields = validFields
	}

	return Model{
		Name:        name,
		LocaleAlpha: localeAlpha,
		Fields:      fields,
		validFields: validFields,
		ProjectID:   projectId,
	}
}

func newView(model declarations.Variable, returnFields []string) map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = model.ID
	m["name"] = model.Name
	m["projectID"] = model.ProjectID

	for _, f := range returnFields {
		if f == "groups" {
			m["groups"] = model.Groups
		}

		if f == "value" {
			m["value"] = model.Value
		}

		if f == "behaviour" {
			m["behaviour"] = model.Behaviour
		}

		if f == "metadata" {
			m["metadata"] = model.Metadata
		}

		if f == "value" {
			m["value"] = model.Value
		}

		if f == "created_at" {
			m["createdAt"] = model.CreatedAt
		}

		if f == "updated_at" {
			m["updatedAt"] = model.UpdatedAt
		}
	}

	return m
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":        a.Name,
		"fieldsValid": a.Fields,
		"locale":      a.LocaleAlpha,
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("name", validation.Required),
			validation.Key("locale", validation.Required, validation.By(func(value interface{}) error {
				t := value.(string)

				if !locales.ExistsByAlpha(t) {
					return errors.New(fmt.Sprintf("Locale '%s' does not exist.", t))
				}

				return nil
			})),
			validation.Key("fieldsValid", validation.By(func(value interface{}) error {
				fields := value.([]string)
				vFields := a.validFields

				for _, f := range fields {
					if !sdk.Includes(vFields, f) {
						return errors.New(fmt.Sprintf("%s is not a valid field to return. Valid fields are %s", f, strings.Join(a.validFields, ", ")))
					}
				}

				return nil
			})),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
