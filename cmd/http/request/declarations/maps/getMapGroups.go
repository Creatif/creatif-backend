package maps

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetMapGroups struct {
	Name      string `param:"name"`
	ProjectID string `param:"projectID"`
}

func SanitizeGetMapGroups(model GetMapGroups) GetMapGroups {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
