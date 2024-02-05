package app

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetGroups struct {
	ProjectID string `param:"projectId"`
}

func SanitizeGetGroups(model GetGroups) GetGroups {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
