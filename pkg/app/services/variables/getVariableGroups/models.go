package getVariableGroups

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type LogicModel struct {
	Group string
}

type Model struct {
	Name      string
	ProjectID string
}

func NewModel(name, projectID string) Model {
	return Model{
		Name:      name,
		ProjectID: projectID,
	}
}

type View struct {
}

func newView(model []LogicModel) []string {
	if len(model) == 0 {

	}
	return sdk.Map(model, func(idx int, value LogicModel) string {
		return value.Group
	})
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"id":        a.Name,
		"projectID": a.ProjectID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("id", validation.Required),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
