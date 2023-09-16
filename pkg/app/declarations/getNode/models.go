package getNode

import (
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"strings"
	"time"
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
	Name   string `json:"name"`
	Fields []string

	validFields []string
}

func NewModel(name string, fields []string) Model {
	if len(fields) == 0 {
		fields = validFields
	}

	return Model{
		Name:        name,
		Fields:      fields,
		validFields: validFields,
	}
}

type Node struct {
	ID string `gorm:"primarykey"`

	Name      string         `gorm:"index;uniqueIndex:unique_node"`
	Behaviour string         // readonly,modifiable
	Groups    pq.StringArray `gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata  datatypes.JSON
	Value     datatypes.JSON

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model Node, returnFields []string) map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = model.ID
	m["name"] = model.Name

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
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("name", validation.Required),
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
