package images

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetImage struct {
	ProjectID   string `param:"projectID"`
	StructureID string `param:"id"`
}

func SanitizeGetImage(model GetImage) GetImage {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.StructureID = p.Sanitize(model.StructureID)

	return model
}
