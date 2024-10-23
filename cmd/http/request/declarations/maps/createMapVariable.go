package maps

import (
	"github.com/microcosm-cc/bluemonday"
)

type CreateMapVariableModel struct {
	Name      string   `json:"name"`
	Metadata  string   `json:"metadata"`
	Locale    string   `json:"locale"`
	Value     string   `json:"value"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
}

type Entry struct {
	Type  string
	Model interface{}
}

type CreateMap struct {
	Variables []CreateMapVariableModel `json:"variables"`
	Name      string                   `json:"name"`
	ProjectID string                   `param:"projectID"`
}

func SanitizeMapModel(model CreateMap) CreateMap {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)

	variables := model.Variables
	newVariables := make([]CreateMapVariableModel, len(variables))
	for i := 0; i < len(model.Variables); i++ {
		m := variables[i]
		m.Name = p.Sanitize(m.Name)
		m.Behaviour = p.Sanitize(m.Behaviour)
		m.Locale = p.Sanitize(m.Locale)

		newGroups := make([]string, len(m.Groups))
		for a := 0; a < len(m.Groups); a++ {
			newGroups[a] = p.Sanitize(m.Groups[a])
		}

		m.Groups = newGroups

		newVariables = append(newVariables, m)
	}

	model.Variables = newVariables

	return model
}
