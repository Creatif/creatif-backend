package declarations

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type AddToMap struct {
	Name      string           `json:"name"`
	ProjectID string           `param:"projectID"`
	Locale    string           `json:"locale"`
	Entry     MapVariableModel `json:"entry"`
}

func SanitizeAddToMap(model AddToMap) AddToMap {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.Name = p.Sanitize(model.Name)

	variable := model.Entry
	variable.Name = p.Sanitize(variable.Name)
	variable.Behaviour = p.Sanitize(variable.Behaviour)
	variable.Groups = sdk.Map(variable.Groups, func(idx int, value string) string {
		return p.Sanitize(value)
	})

	model.Entry = variable

	return model
}
