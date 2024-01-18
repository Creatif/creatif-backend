package maps

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetMap struct {
	Name      string `param:"name"`
	ProjectID string `param:"projectID"`
}

func SanitizeGetMap(model GetMap) GetMap {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
