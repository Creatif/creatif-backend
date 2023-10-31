package variables

import (
	"github.com/microcosm-cc/bluemonday"
)

type DeleteVariable struct {
	Name      string `json:"name"`
	Locale    string `json:"locale"`
	ID        string `json:"id"`
	ShortID   string `json:"shortID"`
	ProjectID string `param:"projectID"`
}

func SanitizeDeleteVariable(model DeleteVariable) DeleteVariable {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.ID = p.Sanitize(model.ID)
	model.ShortID = p.Sanitize(model.ShortID)

	return model
}
