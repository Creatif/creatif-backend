package variables

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetVariableGroups struct {
	Name      string `param:"name"`
	ProjectID string `param:"projectID"`
}

func SanitizeGetVariableGroups(model GetVariableGroups) GetVariableGroups {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
