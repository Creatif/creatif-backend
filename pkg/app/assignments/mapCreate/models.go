package mapCreate

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"time"
)

type AssignValueModel struct {
	Name  string
	Value []byte

	workingMap declarations.Map
}

func NewAssignValueModel(name string, value []byte) AssignValueModel {
	return AssignValueModel{
		Name:  name,
		Value: value,
	}
}

type View struct {
	ID    uuid.UUID   `json:"id"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`

	// createdAt and updatedAt are the values of assignments.MapValueNode, not declarations.MapNode
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LogicModel struct {
	m              declarations.Map
	assignmentNode assignments.MapNode
	valueNode      assignments.MapValueNode
}

func newView(model LogicModel) View {
	return View{
		ID:        model.m.ID,
		Name:      model.m.Name,
		Value:     []byte(model.valueNode.Value),
		CreatedAt: model.assignmentNode.CreatedAt,
		UpdatedAt: model.assignmentNode.UpdatedAt,
	}
}

func (a *AssignValueModel) Validate() map[string]string {
	v := map[string]interface{}{
		"name": a.Name,
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("name", validation.Required, validation.Length(0, 200), validation.By(func(value interface{}) error {
				name := value.(string)
				node := declarations.Map{}
				if err := storage.GetBy(node.TableName(), "name", name, &node); err != nil {
					return errors.New(fmt.Sprintf("Cannot find declaration map with name %s", a.Name))
				}

				a.workingMap = node

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
