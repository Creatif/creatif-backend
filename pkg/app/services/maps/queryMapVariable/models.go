package queryMapVariable

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Model struct {
	Name      string
	ItemID    string
	ProjectID string
}

type QueryReference struct {
	ID                string
	Name              string
	ParentType        string
	ChildType         string
	ParentID          string
	ChildID           string
	ChildStructureID  string
	ParentStructureID string
	StructureName     string
}

func NewModel(projectId, name, itemID string) Model {
	return Model{
		ProjectID: projectId,
		Name:      name,
		ItemID:    itemID,
	}
}

type LogicModel struct {
	Variable  QueryVariable
	Reference []shared.QueryReference
}

type ReferenceView struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	ParentType        string `json:"parentType"`
	ChildType         string `json:"childType"`
	ParentID          string `json:"parentId"`
	ChildID           string `json:"childId"`
	ChildStructureID  string `json:"childStructureId"`
	ParentStructureID string `json:"parentStructureId"`
	StructureName     string `json:"structureName"`
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
		Groups:    model.Variable.Groups,
		Behaviour: model.Variable.Behaviour,
		Metadata:  model.Variable.Metadata,
		Value:     model.Variable.Value,
		CreatedAt: model.Variable.CreatedAt,
		UpdatedAt: model.Variable.UpdatedAt,
		References: sdk.Map(model.Reference, func(idx int, value shared.QueryReference) ReferenceView {
			return ReferenceView{
				ID:                value.ID,
				Name:              value.Name,
				ParentStructureID: value.ParentStructureID,
				ChildStructureID:  value.ChildStructureID,
				ParentID:          value.ParentID,
				ChildID:           value.ChildID,
				ParentType:        value.ParentType,
				ChildType:         value.ChildType,
				StructureName:     value.StructureName,
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
