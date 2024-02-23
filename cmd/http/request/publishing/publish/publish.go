package maps

import (
	"github.com/microcosm-cc/bluemonday"
)

type Publish struct {
	ProjectID string `param:"projectId"`
}

func SanitizePublish(model Publish) Publish {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
