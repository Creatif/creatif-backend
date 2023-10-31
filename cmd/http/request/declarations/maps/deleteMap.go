package maps

import (
	"github.com/microcosm-cc/bluemonday"
)

type DeleteMap struct {
	Name      string `param:"name"`
	ID        string `json:"id"`
	ShortID   string `json:"shortID"`
	ProjectID string `param:"projectID"`
	Locale    string `param:"locale"`
}

func SanitizeDeleteMap(model DeleteMap) DeleteMap {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.ID = p.Sanitize(model.ID)
	model.ShortID = p.Sanitize(model.ShortID)

	return model
}
