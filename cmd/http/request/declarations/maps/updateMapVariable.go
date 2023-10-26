package maps

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type UpdateMapVariable struct {
	ProjectID    string           `param:"projectID"`
	Locale       string           `param:"locale"`
	MapName      string           `param:"mapName"`
	VariableName string           `param:"variableName"`
	Fields       []string         `json:"fields"`
	Entry        MapVariableModel `json:"variable"`

	SanitizedFields []string
}

func SanitizeUpdateMapVariable(model UpdateMapVariable) UpdateMapVariable {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.MapName = p.Sanitize(model.MapName)
	model.VariableName = p.Sanitize(model.VariableName)
	model.Locale = p.Sanitize(model.Locale)

	model.SanitizedFields = sdk.Map(model.Fields, func(idx int, value string) string {
		return p.Sanitize(value)
	})

	variable := model.Entry
	variable.Name = p.Sanitize(variable.Name)
	variable.Behaviour = p.Sanitize(variable.Behaviour)
	variable.Groups = sdk.Map(variable.Groups, func(idx int, value string) string {
		return p.Sanitize(value)
	})

	model.Entry = variable

	return model
}
