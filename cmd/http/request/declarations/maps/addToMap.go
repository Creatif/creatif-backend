package maps

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type AddToMap struct {
	ProjectID string                 `param:"projectID"`
	Variable  CreateMapVariableModel `json:"variable"`
	// Name is the name or ID of the structure
	Name        string       `json:"name"`
	Connections []Connection `json:"connections"`
	ImagePaths  []string     `json:"imagePaths"`
}

type Connection struct {
	Path          string `json:"name"`
	StructureType string `json:"structureType"`
	VariableID    string `json:"variableId"`
}

func SanitizeAddToMap(model AddToMap) AddToMap {
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
			Path:          p.Sanitize(value.Path),
			StructureType: p.Sanitize(value.StructureType),
			VariableID:    p.Sanitize(value.VariableID),
		}
	})

	return model
}
