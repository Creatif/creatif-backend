package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetVariable struct {
	Name      string   `param:"name"`
	Fields    []string `query:"fields"`
	ProjectID string   `param:"projectID"`
}

func SanitizeGetVariable(model GetVariable) GetVariable {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)

	if len(model.Fields) > 0 {
		newFields := make([]string, len(model.Fields))

		for i := 0; i < len(model.Fields); i++ {
			newFields[i] = p.Sanitize(model.Fields[i])
		}

		model.Fields = newFields
	}

	return model
}
