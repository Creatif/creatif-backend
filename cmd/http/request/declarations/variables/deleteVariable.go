package variables

import (
	"github.com/microcosm-cc/bluemonday"
)

type DeleteVariable struct {
	Name      string `json:"name"`
	Locale    string `json:"locale"`
	ProjectID string `param:"projectID"`
}

func SanitizeDeleteVariable(model DeleteVariable) DeleteVariable {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)

	return model
}
