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
	Locale    string   `json:"locale"`
	Value     string   `json:"value"`
}

type CreateList struct {
	Name         string               `json:"name"`
	ProjectID    string               `param:"projectID"`
	Variables    []CreateListVariable `json:"variables"`
	GracefulFail bool                 `json:"gracefulFail"`
}

func SanitizeCreateList(model CreateList) CreateList {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)

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
