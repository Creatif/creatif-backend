package create

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"time"
)

type NodeValidation struct {
	Required    bool
	Length      string
	ExactValue  string
	ExactValues string
	IsDate      bool
}

type CreateNodeModel struct {
	Name       string                    `json:"name"`
	Type       string                    `json:"type"`
	Metadata   []byte                    `json:"metadata"`
	Groups     []string                  `json:"groups"`
	Behaviour  string                    `json:"behaviour"`
	Validation map[string]NodeValidation `json:"validation"`
}

func NewCreateNodeModel(name, t, behaviour string, groups []string, metadata []byte, validation map[string]NodeValidation) CreateNodeModel {
	return CreateNodeModel{
		Name:       name,
		Type:       t,
		Behaviour:  behaviour,
		Groups:     groups,
		Validation: validation,
		Metadata:   metadata,
	}
}

type View struct {
	ID        string                 `json:"name"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Groups    []string               `json:"groups"`
	Behaviour string                 `json:"behaviour"`
	Metadata  map[string]interface{} `json:"metadata"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model declarations.Node) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		Type:      model.Type,
		Groups:    model.Groups,
		Behaviour: model.Behaviour,
		Metadata:  sdk.UnmarshalToMap([]byte(model.Metadata)),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
