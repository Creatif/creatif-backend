package maps

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type UpdateMapVariable struct {
	ProjectID string           `param:"projectID"`
	Locale    string           `param:"locale"`
	Name      string           `json:"name"`
	Entry     MapVariableModel `json:"entry"`
}

func SanitizeUpdateMapVariable(model UpdateMapVariable) UpdateMapVariable {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Name = p.Sanitize(model.Name)
	model.Locale = p.Sanitize(model.Locale)

	variable := model.Entry
	variable.Name = p.Sanitize(variable.Name)
	variable.Behaviour = p.Sanitize(variable.Behaviour)
	variable.Groups = sdk.Map(variable.Groups, func(idx int, value string) string {
		return p.Sanitize(value)
	})

	model.Entry = variable

	return model
}
