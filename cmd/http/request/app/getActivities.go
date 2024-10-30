package app

import "github.com/microcosm-cc/bluemonday"

type GetActivities struct {
	ProjectID string `param:"projectId"`
}

func SanitizeGetActivities(model GetActivities) GetActivities {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
