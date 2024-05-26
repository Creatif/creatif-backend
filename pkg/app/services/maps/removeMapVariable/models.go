package removeMapVariable

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	Name         string
	VariableName string
	ProjectID    string
}

func NewModel(projectId, name, variableName string) Model {
	return Model{
		Name:         name,
		ProjectID:    projectId,
		VariableName: variableName,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID": a.ProjectID,
		"entryName": a.VariableName,
		"name":      a.Name,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required),
			validation.Key("entryName", validation.Required),
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
