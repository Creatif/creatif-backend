package getNode

import (
	"creatif/pkg/lib/sdk"
	"encoding/json"
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

type GetNodeModel struct {
	// this can be project name
	ID     string `json:"id"`
	Fields []string

	validFields []string
}

func NewGetNodeModel(id string, fields []string) GetNodeModel {
	if len(fields) == 0 {
		fields = validFields
	}

	return GetNodeModel{
		ID:          id,
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

func (a *GetNodeModel) Validate() map[string]string {
	v := map[string]interface{}{
		"id":          a.ID,
		"fieldsValid": a.Fields,
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("id", validation.Required),
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
		var e map[string]string
		b, err := json.Marshal(err)
		if err != nil {
			return map[string]string{
				"unrecoverable": "An internal validation error occurred. This should not happen. Please, submit a bug.",
			}
		}

		if err := json.Unmarshal(b, &e); err != nil {
			return map[string]string{
				"unrecoverable": "An internal validation error occurred. This should not happen. Please, submit a bug.",
			}
		}

		return e
	}

	return nil
}
