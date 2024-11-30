package queryMapVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"time"
)

type ConnectionStructureView struct {
	StructureName string `json:"structureName"`
	StructureID   string `json:"structureId"`
	StructureType string `json:"structureType"`
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

type View struct {
	ID        string      `json:"id"`
	Locale    string      `json:"locale"`
	ShortID   string      `json:"shortId"`
	Groups    []string    `json:"groups"`
	Name      string      `json:"name"`
	Behaviour string      `json:"behaviour"`
	Metadata  interface{} `json:"metadata"`
	Value     interface{} `json:"value"`

	Connections     []ConnectionView          `json:"connections"`
	ChildStructures []ConnectionStructureView `json:"childStructures"`

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
		ChildStructures: sdk.Map(model.ChildConnectionStructures, func(idx int, value ChildConnectionStructure) ConnectionStructureView {
			return ConnectionStructureView{
				StructureName: value.StructureName,
				StructureID:   value.StructureID,
				StructureType: value.StructureType,
			}
		}),
	}
}
