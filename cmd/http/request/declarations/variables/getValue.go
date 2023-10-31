package variables

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetValue struct {
	Name      string `param:"name"`
	ShortID   string `json:"shortID"`
	ID        string `json:"id"`
	ProjectID string `param:"projectID"`
	Locale    string `param:"locale"`
}

func SanitizeGetValue(model GetValue) GetValue {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.ID = p.Sanitize(model.ID)
	model.ShortID = p.Sanitize(model.ShortID)

	return model
}
