package maps

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type UpdateMapVariable struct {
	ProjectID string           `param:"projectID"`
	Name      string           `param:"name"`
	ItemID    string           `param:"itemId"`
	Fields    string           `query:"fields"`
	Variable  MapVariableModel `json:"variable"`

	ResolvedFields []string
}

func SanitizeUpdateMapVariable(model UpdateMapVariable) UpdateMapVariable {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Name = p.Sanitize(model.Name)
	model.ItemID = p.Sanitize(model.ItemID)

	model.ResolvedFields = sdk.Map(strings.Split(model.Fields, "|"), func(idx int, value string) string {
		trimmed := strings.Trim(value, " ")
		return p.Sanitize(trimmed)
	})

	variable := model.Variable
	variable.Name = p.Sanitize(variable.Name)
	variable.Locale = p.Sanitize(variable.Locale)
	variable.Behaviour = p.Sanitize(variable.Behaviour)
	variable.Groups = sdk.Map(variable.Groups, func(idx int, value string) string {
		return p.Sanitize(value)
	})

	model.Variable = variable

	return model
}
