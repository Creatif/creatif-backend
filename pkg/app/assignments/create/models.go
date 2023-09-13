package create

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type CreateNodeModel struct {
	Name  string
	Value []byte

	assignedValue   interface{}
	declarationNode declarations.Node
}

func NewCreateNodeModel(name string, value []byte) *CreateNodeModel {
	return &CreateNodeModel{
		Name:  name,
		Value: value,
	}
}

type View struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Value     interface{}    `json:"value"`
	Behaviour string         // readonly,modifiable
	Groups    pq.StringArray `json:"groups"`
	Metadata  datatypes.JSON `json:"metadata"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

func newView(model declarations.Node, value interface{}) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		Value:     value,
		Behaviour: model.Behaviour,
		Groups:    model.Groups,
		Metadata:  model.Metadata,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func (a *CreateNodeModel) Validate() map[string]string {
	v := map[string]interface{}{
		"name":             a.Name,
		"isNodeModifiable": "",
		"value":            a.Value,
		"isValueNull":      a.Value,
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("name", validation.Required, validation.Length(0, 200), validation.By(func(value interface{}) error {
				node := declarations.Node{}
				if err := storage.GetBy(node.TableName(), "name", a.Name, &node); err != nil {
					return errors.New(fmt.Sprintf("Cannot find declaration node with name %s", a.Name))
				}

				a.declarationNode = node

				return nil
			})),
			validation.Key("value", validation.Required),
			validation.Key("isValueNull", validation.By(func(value interface{}) error {
				v := value.([]byte)

				if string(v) == "null" {
					return errors.New("null is not a valid value. Send an empty string if you want to set an empty value.")
				}

				return nil
			})),
			validation.Key("isNodeModifiable", validation.By(func(value interface{}) error {
				if a.declarationNode.ID == "" {
					return nil
				}

				if a.declarationNode.Behaviour != "modifiable" {
					return errors.New("This node is 'readonly' and is not modifiable. You can only assign a value to a 'modifiable' node")
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
