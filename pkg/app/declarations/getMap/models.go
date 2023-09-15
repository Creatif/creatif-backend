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
	Name   string
	Fields []string

	validFields []string
}

func NewModel(name string, fields []string) Model {
	return Model{
		Name:        name,
		Fields:      fields,
		validFields: validFields,
	}
}

type LogicModel struct {
	nodeMap declarations.Map
	nodes   []Node
}

type Node struct {
	ID string `json:"id" gorm:"primarykey"`

	Name      string         `json:"name" gorm:"index;uniqueIndex:unique_node"`
	Value     datatypes.JSON `json:"value"`
	Behaviour string         `json:"behaviour"`
	Groups    pq.StringArray `json:"groups" gorm:"type:text[]"`
	Metadata  datatypes.JSON `json:"metadata"`

	CreatedAt time.Time `json:"createdAt" gorm:"<-:createNode"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type View struct {
	ID    string                   `json:"id"`
	Name  string                   `json:"name"`
	Nodes []map[string]interface{} `json:"nodes"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model LogicModel, returnFields []string) View {
	m := make([]map[string]interface{}, 0)

	for _, n := range model.nodes {
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
		ID:        model.nodeMap.ID,
		Name:      model.nodeMap.Name,
		Nodes:     m,
		CreatedAt: model.nodeMap.CreatedAt,
		UpdatedAt: model.nodeMap.UpdatedAt,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":        a.Name,
		"fieldsValid": a.Fields,
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
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
