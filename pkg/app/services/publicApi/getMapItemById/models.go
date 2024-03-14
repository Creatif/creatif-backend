package getMapItemById

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Options struct {
	ValueOnly bool
}

type Model struct {
	ProjectID string
	ItemID    string

	Options Options
}

func NewModel(projectId, itemId string, options Options) Model {
	return Model{
		ProjectID: projectId,
		ItemID:    itemId,

		Options: options,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID": a.ProjectID,
		"itemId":    a.ItemID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("itemId", validation.Required, validation.RuneLength(26, 26)),
		),
	); err != nil {
		return sdk.ErrorToResponseError(err)
	}

	return nil
}

type ConnectionView struct {
	StructureID      string `json:"structureId"`
	StructureShortID string `json:"structureShortId"`
	StructureName    string `json:"structureName"`
	ConnectionType   string `json:"connectionType"`

	ItemName    string `json:"itemName"`
	ItemID      string `json:"itemId"`
	ItemShortID string `json:"itemShortId"`

	ProjectID string      `json:"projectId"`
	Locale    string      `json:"locale"`
	Index     float64     `json:"index"`
	Groups    []string    `json:"groups"`
	Behaviour string      `json:"behaviour"`
	Value     interface{} `json:"value"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type View struct {
	StructureID      string `json:"structureId,omitempty"`
	StructureShortID string `json:"structureShortId,omitempty"`
	StructureName    string `json:"structureName,omitempty"`

	Name    string `json:"itemName,omitempty"`
	ID      string `json:"itemId,omitempty"`
	ShortID string `json:"itemShortId,omitempty"`

	ProjectID string      `json:"projectId,omitempty"`
	Locale    string      `json:"locale,omitempty"`
	Index     float64     `json:"index,omitempty"`
	Groups    []string    `json:"groups,omitempty"`
	Behaviour string      `json:"behaviour,omitempty"`
	Value     interface{} `json:"value"`

	Connections map[string]ConnectionView `json:"connections,omitempty"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type LogicModel struct {
	Item        Item
	Connections []ConnectionItem
	Options     Options
}

func newView(model LogicModel) interface{} {
	if model.Options.ValueOnly {
		return model.Item.Value
	}

	locale, _ := locales.GetAlphaWithID(model.Item.Locale)
	connections := make(map[string]ConnectionView)
	for _, c := range model.Connections {
		connectionLocale, _ := locales.GetAlphaWithID(model.Item.Locale)

		connections[c.ConnectionName] = ConnectionView{
			StructureID:      c.ID,
			StructureShortID: c.ShortID,
			StructureName:    c.StructureName,
			ConnectionType:   c.ConnectionType,
			ItemName:         c.Name,
			ItemID:           c.ItemID,
			ItemShortID:      c.ItemShortID,
			ProjectID:        c.ProjectID,
			Locale:           connectionLocale,
			Index:            c.Index,
			Groups:           c.Groups,
			Behaviour:        c.Behaviour,
			Value:            c.Value,
			CreatedAt:        c.CreatedAt,
			UpdatedAt:        c.UpdatedAt,
		}
	}

	view := View{
		StructureID:      model.Item.ID,
		StructureShortID: model.Item.ShortID,
		StructureName:    model.Item.StructureName,
		Name:             model.Item.ItemName,
		ID:               model.Item.ItemID,
		ShortID:          model.Item.ItemShortID,
		ProjectID:        model.Item.ProjectID,
		Locale:           locale,
		Index:            model.Item.Index,
		Groups:           model.Item.Groups,
		Behaviour:        model.Item.Behaviour,
		Value:            model.Item.Value,
		Connections:      connections,
		CreatedAt:        nil,
		UpdatedAt:        nil,
	}

	if !model.Options.ValueOnly {
		view.CreatedAt = &model.Item.CreatedAt
		view.UpdatedAt = &model.Item.UpdatedAt
	}

	return view
}
