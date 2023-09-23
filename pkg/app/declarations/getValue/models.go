package getValue

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/datatypes"
)

type Model struct {
	Name      string `json:"name"`
	ProjectID string `json:"projectID"`
}

func NewModel(projectId, name string) Model {
	return Model{Name: name, ProjectID: projectId}
}

type Variable struct {
	Value datatypes.JSON
}

func newView(model Variable) datatypes.JSON {
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
