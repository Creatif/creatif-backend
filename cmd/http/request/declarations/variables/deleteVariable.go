package variables

import (
	"github.com/microcosm-cc/bluemonday"
)

type DeleteVariable struct {
	Name      string `param:"name"`
	Locale    string `param:"locale"`
	ProjectID string `param:"projectID"`
}

func SanitizeDeleteVariable(model DeleteVariable) DeleteVariable {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)

	return model
}
