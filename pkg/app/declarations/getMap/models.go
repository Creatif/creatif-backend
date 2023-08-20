package create

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"time"
)

type GetNodeModel struct {
	// this can be project name or an id of the map
	ID string `json:"id"`
	// TODO: Add project ID prop here
}

func NewGetNodeModel(id string) GetNodeModel {
	return GetNodeModel{
		ID: id,
	}
}

type View struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Groups    []string               `json:"groups"`
	Behaviour string                 `json:"behaviour"`
	Metadata  map[string]interface{} `json:"metadata"`
	Value     interface{}            `json:"value"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(models []declarations.Node) map[string]View {
	view := make(map[string]View)

	for _, model := range models {
		view[model.Name] = View{
			ID: model.ID,
			// TODO: add value later
			Name:      model.Name,
			Type:      model.Type,
			Groups:    model.Groups,
			Behaviour: model.Behaviour,
			Metadata:  sdk.UnmarshalToMap([]byte(model.Metadata)),
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		}
	}

	return view
}
