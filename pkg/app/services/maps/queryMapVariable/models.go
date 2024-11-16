package queryMapVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Model struct {
	Name      string
	ItemID    string
	ProjectID string
}

type ConnectionView struct {
	ProjectID string `json:"projectId"`

	Path                string `json:"path"`
	ParentVariableID    string `json:"parentVariableId"`
	ParentStructureType string `json:"parentStructureType"`

	ChildVariableID    string `json:"childVariableId"`
	ChildStructureType string `json:"childStructureType"`

	CreatedAt time.Time `json:"createdAt"`
}

func NewModel(projectId, name, itemID string) Model {
	return Model{
		ProjectID: projectId,
		Name:      name,
		ItemID:    itemID,
	}
}

type LogicModel struct {
	Variable    QueryVariable
	Connections []declarations.Connection
}

type View struct {
	ID        string      `json:"id"`
	Locale    string      `json:"locale"`
	ShortID   string      `json:"shortId"`
	Groups    []string    `json:"groups"`
	Name      string      `json:"name"`
	Behaviour string      `json:"behaviour"`
	Metadata  interface{} `json:"metadata"`
	Value     interface{} `json:"value"`

	Connections []ConnectionView `json:"connections"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model LogicModel) View {
	alpha, _ := locales.GetAlphaWithID(model.Variable.LocaleID)
	return View{
		ID:        model.Variable.ID,
		Locale:    alpha,
		ShortID:   model.Variable.ShortID,
		Name:      model.Variable.Name,
		Groups:    model.Variable.Groups,
		Behaviour: model.Variable.Behaviour,
		Metadata:  model.Variable.Metadata,
		Value:     model.Variable.Value,
		CreatedAt: model.Variable.CreatedAt,
		UpdatedAt: model.Variable.UpdatedAt,
		Connections: sdk.Map(model.Connections, func(idx int, value declarations.Connection) ConnectionView {
			return ConnectionView{
				Path:                value.Path,
				ParentVariableID:    value.ParentVariableID,
				ParentStructureType: value.ParentStructureType,
				ChildVariableID:     value.ChildVariableID,
				ChildStructureType:  value.ChildStructureType,
				CreatedAt:           value.CreatedAt,
			}
		}),
	}
}

func (a *Model) Validate() map[string]string {
	v := map[string]interface{}{
		"name":      a.Name,
		"itemId":    a.ItemID,
		"projectID": a.ProjectID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("name", validation.Required),
			validation.Key("itemId", validation.Required),
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
