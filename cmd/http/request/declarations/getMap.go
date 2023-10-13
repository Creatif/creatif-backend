package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetMap struct {
	Name      string `param:"name"`
	Fields    string `query:"fields"`
	ProjectID string `param:"projectID"`
	Locale    string `param:"locale"`
}

func SanitizeGetMap(model GetMap) GetMap {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Fields = p.Sanitize(model.Fields)
	model.Locale = p.Sanitize(model.Locale)

	return model
}
