package getValue

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/datatypes"
)

type Model struct {
	// this can be project name
	Name string `json:"name"`
}

func NewModel(name string) Model {
	return Model{Name: name}
}

type Node struct {
	Value datatypes.JSON
}

func newView(model Node) datatypes.JSON {
	return model.Value
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name": a.Name,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
