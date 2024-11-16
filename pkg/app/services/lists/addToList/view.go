package addToList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"github.com/lib/pq"
	"time"
)

type ConnectionView struct {
	ProjectID string `json:"projectId"`

	Path                string `json:"path"`
	ParentVariableID    string `json:"parentVariableId"`
	ParentStructureType string `json:"parentStructureType"`

	ChildVariableID    string `json:"childVariableId"`
	ChildStructureType string `json:"childStructureType"`

	CreatedAt time.Time `json:"createdAt"`
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

	Connections []ConnectionView `json:"connections"`
}

func newView(model LogicModel) View {
	var m interface{} = model.Variable.Metadata
	if len(model.Variable.Metadata) == 0 {
		m = nil
	}

	var v interface{} = model.Variable.Value
	if len(model.Variable.Value) == 0 {
		v = nil
	}

	return View{
		ID:        model.Variable.ID,
		ShortID:   model.Variable.ShortID,
		Index:     model.Variable.Index,
		Name:      model.Variable.Name,
		Behaviour: model.Variable.Behaviour,
		Groups:    model.Groups,
		Metadata:  m,
		Value:     v,
		Locale:    model.Variable.LocaleID,
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
