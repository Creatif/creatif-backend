package addToMap

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"github.com/lib/pq"
	"time"
)

type ReferenceView struct {
	ID string `json:"id"`

	ParentType string `json:"parentType"`
	ChildType  string `json:"childType"`

	// must be structure type item
	ParentID      string `json:"parentId"`
	ParentShortID string `json:"parentShortId"`
	// must be entire structure
	ChildID      string `json:"childId"`
	ChildShortID string `json:"childShortId"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type View struct {
	ID      string  `json:"id"`
	ShortID string  `json:"shortId"`
	Index   float64 `json:"index"`

	Name      string         `json:"name"`
	Behaviour string         `json:"behaviour"`
	Groups    pq.StringArray `json:"groups"`
	Metadata  interface{}    `json:"metadata"`
	Value     interface{}    `json:"value"`

	Locale string `json:"locale"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	References []ReferenceView `json:"references"`
}

func newView(model LogicModel) View {
	return View{
		ID:        model.Variable.ID,
		ShortID:   model.Variable.ShortID,
		Index:     model.Variable.Index,
		Name:      model.Variable.Name,
		Behaviour: model.Variable.Behaviour,
		Groups:    model.Variable.Groups,
		Metadata:  model.Variable.Metadata,
		Value:     model.Variable.Value,
		Locale:    model.Variable.LocaleID,
		CreatedAt: model.Variable.CreatedAt,
		UpdatedAt: model.Variable.UpdatedAt,
		References: sdk.Map(model.References, func(idx int, value declarations.Reference) ReferenceView {
			return ReferenceView{
				ID:            value.ID,
				ParentType:    value.ParentType,
				ChildType:     value.ChildType,
				ParentID:      value.ParentID,
				ParentShortID: value.ParentShortID,
				ChildID:       value.ChildID,
				ChildShortID:  value.ChildShortID,
				CreatedAt:     value.CreatedAt,
				UpdatedAt:     value.UpdatedAt,
			}
		}),
	}
}
