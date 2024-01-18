package maps

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type UpdateMapVariable struct {
	ProjectID  string            `param:"projectID"`
	Name       string            `param:"name"`
	ItemID     string            `param:"itemId"`
	Fields     string            `query:"fields"`
	Variable   MapVariableModel  `json:"variable"`
	References []UpdateReference `json:"reference"`

	ResolvedFields []string
}

type UpdateReference struct {
	ID            string
	StructureName string
	StructureType string
	VariableID    string
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

	if len(model.References) != 0 {
		model.References = sdk.Map(model.References, func(idx int, value UpdateReference) UpdateReference {
			return UpdateReference{
				ID:            p.Sanitize(value.ID),
				StructureName: p.Sanitize(value.StructureName),
				StructureType: p.Sanitize(value.StructureType),
				VariableID:    p.Sanitize(value.VariableID),
			}
		})
	}

	return model
}
