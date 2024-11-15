package lists

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type AddToList struct {
	ProjectID   string        `param:"projectID"`
	Variable    VariableModel `json:"variable"`
	Name        string        `json:"name"`
	Connections []Connection  `json:"connections"`
	ImagePaths  []string      `json:"imagePaths"`
}

type VariableModel struct {
	Name      string   `json:"name"`
	Metadata  string   `json:"metadata"`
	Locale    string   `json:"locale"`
	Value     string   `json:"value"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
}

type Connection struct {
	Name          string `json:"name"`
	StructureType string `json:"structureType"`
	VariableID    string `json:"variableId"`
}

func SanitizeAddToList(model AddToList) AddToList {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Name = p.Sanitize(model.Name)

	variable := model.Variable
	variable.Name = p.Sanitize(variable.Name)
	variable.Behaviour = p.Sanitize(variable.Behaviour)
	variable.Locale = p.Sanitize(variable.Locale)
	variable.Groups = sdk.Map(variable.Groups, func(idx int, value string) string {
		return p.Sanitize(value)
	})
	model.ImagePaths = sdk.Map(model.ImagePaths, func(idx int, value string) string {
		return p.Sanitize(value)
	})

	model.Variable = variable

	model.Connections = sdk.Map(model.Connections, func(idx int, value Connection) Connection {
		return Connection{
			Name:          p.Sanitize(value.Name),
			StructureType: p.Sanitize(value.StructureType),
			VariableID:    p.Sanitize(value.VariableID),
		}
	})

	return model
}
