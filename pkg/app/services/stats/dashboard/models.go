package dashboard

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	ProjectID string
}

type LogicModel struct {
	StructureID string
	Name        string
	Count       string
	Type        string
	CreatedAt   string
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
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}

type View struct {
	StructureID string `json:"structureId"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Count       string `json:"count"`
	CreatedAt   string `json:"createdAt"`
}

func newView(model []LogicModel) []View {
	return sdk.Map(model, func(idx int, value LogicModel) View {
		return View{
			StructureID: value.StructureID,
			Name:        value.Name,
			Count:       value.Count,
			Type:        value.Type,
			CreatedAt:   value.CreatedAt,
		}
	})
}
