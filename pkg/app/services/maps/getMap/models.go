package getMap

import (
	"creatif/pkg/app/domain/declarations"
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

type Model struct {
	Name      string
	Fields    []string
	ProjectID string
	Groups    []string

	validFields []string
}

func NewModel(projectId, name string, fields []string, groups []string) Model {
	return Model{
		Name:        name,
		ProjectID:   projectId,
		Fields:      fields,
		Groups:      groups,
		validFields: validFields,
	}
}

type LogicModel struct {
	variableMap declarations.Map
	variables   []Variable
}

type Variable struct {
	ID string `json:"id" gorm:"primarykey"`

	Name      string         `json:"name" gorm:"index;uniqueIndex:unique_variable"`
	Value     datatypes.JSON `json:"value"`
	Behaviour string         `json:"behaviour"`
	Groups    pq.StringArray `json:"groups" gorm:"type:text[]"`
	Metadata  datatypes.JSON `json:"metadata"`

	CreatedAt time.Time `json:"createdAt" gorm:"<-:createProject"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type View struct {
	ID        string                   `json:"id"`
	Name      string                   `json:"name"`
	ProjectID string                   `json:"projectID"`
	Variables []map[string]interface{} `json:"variables"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model LogicModel, returnFields []string) View {
	m := make([]map[string]interface{}, 0)

	for _, n := range model.variables {
		o := make(map[string]interface{})

		o["id"] = n.ID
		o["name"] = n.Name

		for _, f := range returnFields {
			if f == "groups" {
				o["groups"] = n.Groups
			}

			if f == "value" {
				o["value"] = n.Value
			}

			if f == "behaviour" {
				o["behaviour"] = n.Behaviour
			}

			if f == "metadata" {
				o["metadata"] = n.Metadata
			}

			if f == "value" {
				o["value"] = n.Value
			}

			if f == "created_at" {
				o["createdAt"] = n.CreatedAt
			}

			if f == "updated_at" {
				o["updatedAt"] = n.UpdatedAt
			}
		}

		m = append(m, o)
	}

	return View{
		ID:        model.variableMap.ID,
		ProjectID: model.variableMap.ProjectID,
		Name:      model.variableMap.Name,
		Variables: m,
		CreatedAt: model.variableMap.CreatedAt,
		UpdatedAt: model.variableMap.UpdatedAt,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":        a.Name,
		"projectID":   a.ProjectID,
		"fieldsValid": a.Fields,
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("name", validation.When(a.Name != "", validation.RuneLength(1, 200))),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("fieldsValid", validation.By(func(value interface{}) error {
				fields := value.([]string)
				vFields := a.validFields

				if len(fields) > 0 {
					for _, f := range fields {
						if !sdk.Includes(vFields, f) {
							return errors.New(fmt.Sprintf("%s is not a valid field to return. Valid fields are %s", f, strings.Join(a.validFields, ", ")))
						}
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
