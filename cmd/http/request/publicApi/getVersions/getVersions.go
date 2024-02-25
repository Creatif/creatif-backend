package getVersions

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetVersions struct {
	ProjectID string `param:"projectId"`
}

func SanitizeGetVersions(model GetVersions) GetVersions {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
