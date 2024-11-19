package maps

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type UpdateMapVariableModel struct {
	Name      string   `json:"name"`
	Metadata  string   `json:"metadata"`
	Locale    string   `json:"locale"`
	Value     string   `json:"value"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
}

type UpdateMapVariable struct {
	ProjectID   string                 `param:"projectID"`
	Name        string                 `param:"name"`
	ItemID      string                 `param:"itemId"`
	Fields      string                 `query:"fields"`
	Variable    UpdateMapVariableModel `json:"variable"`
	Connections []UpdateConnection     `json:"connections"`
	ImagePaths  []string               `json:"imagePaths"`

	ResolvedFields []string
}

type UpdateConnection struct {
	Name          string `json:"name"`
	StructureType string `json:"structureType"`
	VariableID    string `json:"variableId"`
}

func SanitizeUpdateMapVariable(model UpdateMapVariable) UpdateMapVariable {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Name = p.Sanitize(model.Name)
	model.ItemID = p.Sanitize(model.ItemID)

	model.ImagePaths = sdk.Map(model.ImagePaths, func(idx int, value string) string {
		return p.Sanitize(value)
	})

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

	if len(model.Connections) != 0 {
		model.Connections = sdk.Map(model.Connections, func(idx int, value UpdateConnection) UpdateConnection {
			return UpdateConnection{
				Name:          p.Sanitize(value.Name),
				StructureType: p.Sanitize(value.StructureType),
				VariableID:    p.Sanitize(value.VariableID),
			}
		})
	}

	return model
}
