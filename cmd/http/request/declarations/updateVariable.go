package declarations

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type UpdateVariableValues struct {
	Name      string   `json:"name"`
	Metadata  string   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Value     string   `json:"value"`
}

type UpdateVariable struct {
	Fields    []string             `json:"fields"`
	Name      string               `json:"name"`
	Values    UpdateVariableValues `json:"values"`
	ProjectID string               `param:"projectID"`
}

func SanitizeUpdateVariable(model UpdateVariable) UpdateVariable {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)

	newFields := sdk.Map(model.Fields, func(idx int, value string) string {
		return p.Sanitize(value)
	})

	values := UpdateVariableValues{
		Name:     p.Sanitize(model.Values.Name),
		Metadata: model.Values.Metadata,
		Groups: sdk.Map(model.Values.Groups, func(idx int, value string) string {
			return p.Sanitize(value)
		}),
		Behaviour: p.Sanitize(model.Values.Behaviour),
		Value:     model.Values.Value,
	}

	model.Fields = newFields
	model.Values = values

	return model
}
