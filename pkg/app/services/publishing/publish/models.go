package publish

import (
	"creatif/pkg/app/domain/published"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	ProjectID string
	Name      string
}

func NewModel(projectId, name string) Model {
	return Model{
		ProjectID: projectId,
		Name:      name,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID": a.ProjectID,
		"name":      a.Name,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}

type View struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func newView(model published.Version) View {
	return View{
		ID:   model.ID,
		Name: model.Name,
	}
}
