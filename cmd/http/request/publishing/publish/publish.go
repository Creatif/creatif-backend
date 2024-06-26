package publish

import (
	"github.com/microcosm-cc/bluemonday"
)

type Publish struct {
	ProjectID string `param:"projectId"`
	Name      string `json:"name"`
}

func SanitizePublish(model Publish) Publish {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Name = p.Sanitize(model.Name)

	return model
}
