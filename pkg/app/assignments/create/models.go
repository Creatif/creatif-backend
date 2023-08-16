package create

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateNodeModel struct {
	Name  string
	Value interface{}
	Type  string

	declarationNode declarations.Node
}

type AssignNodeTextModel struct {
	Name  string
	Value string
}

type AssignNodeBooleanModel struct {
	Name  string
	Value bool
}

func NewCreateNodeModel(name, t string, value interface{}) *CreateNodeModel {
	return &CreateNodeModel{
		Name:  name,
		Type:  t,
		Value: value,
	}
}

type View struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func newView(model assignments.Node) View {
	return View{
		ID:   model.ID,
		Name: model.Name,
	}
}

func (a *CreateNodeModel) Validate() map[string]string {
	v := map[string]interface{}{
		"name": a.Name,
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
		),
	); err != nil {
		var e map[string]string
		b, _ := json.Marshal(err)
		json.Unmarshal(b, &e)

		return e
	}

	return nil
}
