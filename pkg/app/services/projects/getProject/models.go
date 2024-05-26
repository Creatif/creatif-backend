package getProject

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Model struct {
	ProjectID string
}

func NewModel(projectID string) Model {
	return Model{
		ProjectID: projectID,
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID": a.ProjectID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(1, 200)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}

type View struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	APIKey string `json:"apiKey"`
	Secret string `json:"secret"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model app.Project) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
