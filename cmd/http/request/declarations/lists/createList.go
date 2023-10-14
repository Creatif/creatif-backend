package lists

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type CreateListVariable struct {
	Name      string   `json:"name"`
	Metadata  string   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Value     string   `json:"value"`
}

type CreateList struct {
	Name      string               `json:"name"`
	ProjectID string               `param:"projectID"`
	Locale    string               `json:"locale"`
	Variables []CreateListVariable `json:"variables"`
}

func SanitizeCreateList(model CreateList) CreateList {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)

	newVariables := make([]CreateListVariable, len(model.Variables))
	for i, variable := range model.Variables {
		newGroups := sdk.Sanitize(variable.Groups, func(k int, v string) string {
			return p.Sanitize(v)
		})

		newVariable := CreateListVariable{
			Name:      p.Sanitize(variable.Name),
			Metadata:  variable.Metadata,
			Groups:    newGroups,
			Behaviour: p.Sanitize(variable.Behaviour),
			Value:     variable.Value,
		}

		newVariables[i] = newVariable
	}

	model.Variables = newVariables

	return model
}
