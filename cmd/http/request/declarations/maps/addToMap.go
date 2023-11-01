package maps

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type AddToMap struct {
	Locale    string           `param:"locale"`
	ProjectID string           `param:"projectID"`
	Entry     MapVariableModel `json:"entry"`
	Name      string           `json:"name"`
	ID        string           `json:"id"`
	ShortID   string           `json:"shortID"`
}

func SanitizeAddToMap(model AddToMap) AddToMap {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.Name = p.Sanitize(model.Name)
	model.ID = p.Sanitize(model.ID)
	model.ShortID = p.Sanitize(model.ShortID)

	variable := model.Entry
	variable.Name = p.Sanitize(variable.Name)
	variable.Behaviour = p.Sanitize(variable.Behaviour)
	variable.Groups = sdk.Map(variable.Groups, func(idx int, value string) string {
		return p.Sanitize(value)
	})

	model.Entry = variable

	return model
}
