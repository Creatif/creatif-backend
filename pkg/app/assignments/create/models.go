package create

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
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
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Value     interface{}    `json:"value"`
	Type      string         // text,image,file,boolean
	Behaviour string         // readonly,modifiable
	Groups    pq.StringArray `json:"groups"`
	Metadata  datatypes.JSON `json:"metadata"`
}

func newView(model declarations.Node, value interface{}) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		Value:     value,
		Type:      model.Type,
		Behaviour: model.Behaviour,
		Groups:    model.Groups,
		Metadata:  model.Metadata,
	}
}

func (a *CreateNodeModel) Validate() map[string]string {
	v := map[string]interface{}{
		"name":             a.Name,
		"isNodeModifiable": "",
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
			validation.Key("isNodeModifiable", validation.By(func(value interface{}) error {
				if a.declarationNode.ID.String() == "" {
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
