package updatePublished

import (
	"github.com/microcosm-cc/bluemonday"
)

type UpdatePublished struct {
	ProjectID string `param:"projectId"`
	Name      string `json:"name"`
}

func SanitizeUpdatePublished(model UpdatePublished) UpdatePublished {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Name = p.Sanitize(model.Name)

	return model
}
