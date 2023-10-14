package maps

import (
	"github.com/microcosm-cc/bluemonday"
)

type DeleteMap struct {
	Name      string `param:"name"`
	ProjectID string `param:"projectID"`
	Locale    string `param:"locale"`
}

func SanitizeDeleteMap(model DeleteMap) DeleteMap {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)

	return model
}
