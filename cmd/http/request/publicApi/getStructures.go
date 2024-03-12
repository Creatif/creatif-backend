package publicApi

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetStructures struct {
	ProjectID string `param:"projectId"`
}

func SanitizeGetStructures(model GetStructures) GetStructures {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
