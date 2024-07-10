package files

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetFile struct {
	ProjectID   string `param:"projectID"`
	StructureID string `param:"id"`
}

func SanitizeGetFile(model GetFile) GetFile {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.StructureID = p.Sanitize(model.StructureID)

	return model
}
