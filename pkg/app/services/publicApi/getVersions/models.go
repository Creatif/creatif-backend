package getVersions

import (
	"creatif/pkg/app/domain/published"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
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
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model []published.Version) []View {
	return sdk.Map(model, func(idx int, value published.Version) View {
		return View{
			ID:        value.ID,
			Name:      value.Name,
			ProjectID: value.ProjectID,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}
	})
}
