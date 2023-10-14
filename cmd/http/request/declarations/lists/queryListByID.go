package lists

import (
	"github.com/microcosm-cc/bluemonday"
)

type QueryListByID struct {
	Name      string `param:"name"`
	ID        string `param:"id"`
	ProjectID string `param:"projectID"`
	Locale    string `param:"locale"`
}

func SanitizeQueryListByID(model QueryListByID) QueryListByID {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.ID = p.Sanitize(model.ID)

	return model
}
