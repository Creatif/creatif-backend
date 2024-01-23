package lists

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type AddToList struct {
	ProjectID  string        `param:"projectID"`
	Variable   VariableModel `json:"variable"`
	Name       string        `json:"name"`
	References []Reference   `json:"references"`
}

type VariableModel struct {
	Name      string   `json:"name"`
	Metadata  string   `json:"metadata"`
	Locale    string   `json:"locale"`
	Value     string   `json:"value"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
}

type Reference struct {
	Name          string `json:"name"`
	StructureName string `json:"structureName"`
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

	model.Variable = variable

	model.References = sdk.Map(model.References, func(idx int, value Reference) Reference {
		return Reference{
			Name:          p.Sanitize(value.Name),
			StructureName: p.Sanitize(value.StructureName),
			StructureType: p.Sanitize(value.StructureType),
			VariableID:    p.Sanitize(value.VariableID),
		}
	})

	return model
}
