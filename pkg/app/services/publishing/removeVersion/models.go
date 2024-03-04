package removeVersion

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	ProjectID string
	ID        string
}

func NewModel(projectId, id string) Model {
	return Model{
		ProjectID: projectId,
		ID:        id,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID": a.ProjectID,
		"id":        a.ID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("id", validation.Required, validation.RuneLength(26, 26)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
