package getValue

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/datatypes"
)

type Model struct {
	// this can be project name
	ID string `json:"id"`
}

func NewModel(id string) Model {
	return Model{ID: id}
}

type Node struct {
	Value datatypes.JSON
}

func newView(model Node) datatypes.JSON {
	return model.Value
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"id": a.ID,
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("id", validation.Required),
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
