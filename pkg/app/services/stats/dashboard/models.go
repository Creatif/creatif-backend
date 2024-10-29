package dashboard

import (
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Model struct {
	ProjectID string
}

type StructureLogicModel struct {
	ID        string
	Name      string
	Count     int
	Type      string
	CreatedAt string
}

type VersionLogicModel struct {
	ID        string
	Name      string
	CreatedAt string
}

type LogicModel struct {
	Structures []StructureLogicModel
	Versions   []VersionLogicModel
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

type StructureView struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Count     int    `json:"count"`
	CreatedAt string `json:"createdAt"`
}

type VersionView struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type View struct {
	Structures []StructureView `json:"structures"`
	Versions   []VersionView   `json:"versions"`
}

func newView(model LogicModel) View {
	return View{
		Structures: sdk.Map(model.Structures, func(idx int, value StructureLogicModel) StructureView {
			return StructureView{
				ID:        value.ID,
				Name:      value.Name,
				Type:      value.Type,
				Count:     value.Count,
				CreatedAt: value.CreatedAt,
			}
		}),
		Versions: sdk.Map(model.Versions, func(idx int, value VersionLogicModel) VersionView {
			return VersionView{
				ID:        value.ID,
				Name:      value.Name,
				CreatedAt: value.CreatedAt,
			}
		}),
	}
}
