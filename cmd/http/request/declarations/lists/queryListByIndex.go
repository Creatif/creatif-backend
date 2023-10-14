package lists

import (
	"github.com/microcosm-cc/bluemonday"
)

type QueryListByIndex struct {
	Name      string `param:"name"`
	Index     int64  `param:"index"`
	ProjectID string `param:"projectID"`
	Locale    string `param:"locale"`
}

func SanitizeQueryListByIndex(model QueryListByIndex) QueryListByIndex {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)

	return model
}
