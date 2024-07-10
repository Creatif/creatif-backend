package getFile

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	ProjectID   string
	StructureID string
}

func NewModel(projectId, structureId string) Model {
	return Model{
		ProjectID:   projectId,
		StructureID: structureId,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID":   a.ProjectID,
		"structureID": a.StructureID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("structureID", validation.Required, validation.RuneLength(27, 27)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
