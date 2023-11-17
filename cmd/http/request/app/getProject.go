package app

import "github.com/microcosm-cc/bluemonday"

type GetProject struct {
	ProjectID string `param:"id"`
}

func SanitizeGetProject(model GetProject) GetProject {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
