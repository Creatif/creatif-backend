package getFile

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	ProjectID string
	Version   string
	FileID    string
}

func NewModel(projectId, fileId, version string) Model {
	return Model{
		ProjectID: projectId,
		Version:   version,
		FileID:    fileId,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID": a.ProjectID,
		"version":   a.Version,
		"fileId":    a.FileID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("fileId", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("version", validation.Required),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
