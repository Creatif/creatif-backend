package removeVersion

import (
	"github.com/microcosm-cc/bluemonday"
)

type RemoveVersion struct {
	ProjectID string `param:"projectId"`
	ID        string `param:"id"`
}

func SanitizeRemoveVersion(model RemoveVersion) RemoveVersion {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ID = p.Sanitize(model.ID)

	return model
}
