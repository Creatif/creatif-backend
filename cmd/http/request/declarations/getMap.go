package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetMap struct {
	Name      string `param:"name"`
	Fields    string `query:"fields"`
	ProjectID string `param:"projectID"`
}

func SanitizeGetMap(model GetMap) GetMap {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.Fields = p.Sanitize(model.Fields)

	return model
}
