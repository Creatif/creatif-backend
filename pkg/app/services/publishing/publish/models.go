package publish

import (
	"creatif/pkg/app/domain/published"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	ProjectID string
}

func NewModel(projectId string) Model {
	return Model{
		ProjectID: projectId,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID": a.ProjectID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}

type View struct {
	VersionID string `json:"versionId"`
}

func newView(model published.Version) View {
	return View{
		VersionID: model.ID,
	}
}
