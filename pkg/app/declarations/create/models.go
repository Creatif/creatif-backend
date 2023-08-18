package create

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"time"
)

type ValidationLength struct {
	Min   int `json:"min"`
	Max   int `json:"max"`
	Exact int `json:"exact"`
}

type NodeValidation struct {
	Required    bool
	Length      ValidationLength
	ExactValue  string
	ExactValues []string
	IsDate      bool
}

type CreateNodeModel struct {
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Metadata   []byte         `json:"metadata"`
	Groups     []string       `json:"groups"`
	Behaviour  string         `json:"behaviour"`
	Validation NodeValidation `json:"validation"`
}

func NewCreateNodeModel(name, t, behaviour string, groups []string, metadata []byte, validation NodeValidation) CreateNodeModel {
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
	ID        string                 `json:"id"`
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
