package getStructures

import (
	"creatif/pkg/app/domain/declarations"
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
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}

type View struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ShortID       string `json:"shortId"`
	StructureType string `json:"structureType"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LogicModel struct {
	Lists []declarations.List `json:"lists"`
	Maps  []declarations.Map  `json:"maps"`
}

func newView(model LogicModel) []View {
	lists := sdk.Map(model.Lists, func(idx int, value declarations.List) View {
		return View{
			ID:            value.ID,
			Name:          value.Name,
			ShortID:       value.ShortID,
			StructureType: "list",
			CreatedAt:     value.CreatedAt,
			UpdatedAt:     value.UpdatedAt,
		}
	})

	maps := sdk.Map(model.Maps, func(idx int, value declarations.Map) View {
		return View{
			ID:            value.ID,
			Name:          value.Name,
			ShortID:       value.ShortID,
			StructureType: "map",
			CreatedAt:     value.CreatedAt,
			UpdatedAt:     value.UpdatedAt,
		}
	})

	views := make([]View, 0)
	for _, l := range lists {
		views = append(views, l)
	}

	for _, m := range maps {
		views = append(views, m)
	}

	return views
}
