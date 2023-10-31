package variables

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetVariable struct {
	Name      string   `param:"name"`
	ShortID   string   `json:"shortID"`
	ID        string   `json:"id"`
	Fields    []string `query:"fields"`
	ProjectID string   `param:"projectID"`
	Locale    string   `param:"locale"`
}

func SanitizeGetVariable(model GetVariable) GetVariable {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.ID = p.Sanitize(model.ID)
	model.ShortID = p.Sanitize(model.ShortID)

	if len(model.Fields) > 0 {
		newFields := make([]string, len(model.Fields))

		for i := 0; i < len(model.Fields); i++ {
			newFields[i] = p.Sanitize(model.Fields[i])
		}

		model.Fields = newFields
	}

	return model
}
