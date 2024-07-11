package publicApi

import "github.com/microcosm-cc/bluemonday"

type GetFile struct {
	ProjectID string `param:"projectID"`
	FileID    string `param:"id"`
	Version   string `param:"version"`
}

func SanitizeGetFile(model GetFile) GetFile {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Version = p.Sanitize(model.Version)
	model.FileID = p.Sanitize(model.FileID)

	return model
}
