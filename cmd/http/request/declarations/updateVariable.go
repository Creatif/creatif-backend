package declarations

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type UpdateVariableValues struct {
	Name      string   `json:"name"`
	Metadata  []byte   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Value     []byte   `json:"value"`
}

type UpdateVariable struct {
	Fields []string             `json:"fields"`
	Name   string               `json:"name"`
	Values UpdateVariableValues `json:"values"`
}

func SanitizeUpdateVariable(model UpdateVariable) UpdateVariable {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)

	newFields := sdk.Map(model.Fields, func(idx int, value string) string {
		return p.Sanitize(value)
	})

	values := UpdateVariableValues{
		Name:     p.Sanitize(model.Values.Name),
		Metadata: []byte(p.Sanitize(string(model.Values.Metadata))),
		Groups: sdk.Map(model.Values.Groups, func(idx int, value string) string {
			return p.Sanitize(value)
		}),
		Behaviour: p.Sanitize(model.Values.Behaviour),
		Value:     []byte(p.Sanitize(string(model.Values.Metadata))),
	}

	model.Fields = newFields
	model.Values = values

	return model
}
