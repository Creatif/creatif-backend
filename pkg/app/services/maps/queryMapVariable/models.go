package queryMapVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lib/pq"
	"time"
)

type Model struct {
	Name      string
	ItemID    string
	ProjectID string
}

func NewModel(projectId, name, itemID string) Model {
	return Model{
		ProjectID: projectId,
		Name:      name,
		ItemID:    itemID,
	}
}

type LogicModel struct {
	Variable  declarations.MapVariable
	Reference []declarations.Reference
}

type ReferenceView struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	StructureID   string `json:"structureId"`
	OwnerID       string `json:"ownerId"`
	StructureType string `json:"structureType"`
}

type View struct {
	ID         string          `json:"id"`
	Locale     string          `json:"locale"`
	ShortID    string          `json:"shortId"`
	Name       string          `json:"name"`
	Behaviour  string          `json:"behaviour"`
	Groups     pq.StringArray  `json:"groups"`
	Metadata   interface{}     `json:"metadata"`
	Value      interface{}     `json:"value"`
	References []ReferenceView `json:"references"`

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
		Behaviour: model.Variable.Behaviour,
		Groups:    model.Variable.Groups,
		Metadata:  model.Variable.Metadata,
		Value:     model.Variable.Value,
		CreatedAt: model.Variable.CreatedAt,
		UpdatedAt: model.Variable.UpdatedAt,
		References: sdk.Map(model.Reference, func(idx int, value declarations.Reference) ReferenceView {
			return ReferenceView{
				ID:            value.ID,
				Name:          value.Name,
				StructureID:   value.ParentID,
				OwnerID:       value.ChildID,
				StructureType: value.ParentType,
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
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}
