package getListItemById

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Model struct {
	ProjectID   string
	ItemID      string
	VersionName string
}

func NewModel(projectId, itemId, versionName string) Model {
	return Model{
		ProjectID:   projectId,
		ItemID:      itemId,
		VersionName: versionName,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID":   a.ProjectID,
		"versionName": a.VersionName,
		"itemId":      a.ItemID,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("itemId", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("versionName", validation.Required),
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
	StructureID      string `json:"structureId"`
	StructureShortID string `json:"structureShortId"`
	StructureName    string `json:"structureName"`

	ItemName    string `json:"itemName"`
	ItemID      string `json:"itemId"`
	ItemShortID string `json:"itemShortId"`

	ProjectID string      `json:"projectId"`
	Locale    string      `json:"locale"`
	Index     float64     `json:"index"`
	Groups    []string    `json:"groups"`
	Behaviour string      `json:"behaviour"`
	Value     interface{} `json:"value"`

	Connections map[string]ConnectionView `json:"connections"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LogicModel struct {
	Item        MapItem
	Connections []ConnectionMapItem
}

func newView(model LogicModel) View {
	locale, _ := locales.GetAlphaWithID(model.Item.Locale)
	connections := make(map[string]ConnectionView)
	for _, c := range model.Connections {
		connectionLocale, _ := locales.GetAlphaWithID(model.Item.Locale)

		connections[c.ConnectionName] = ConnectionView{
			StructureID:      c.ID,
			StructureShortID: c.ShortID,
			StructureName:    c.StructureName,
			ConnectionType:   c.ConnectionType,
			ItemName:         c.ItemName,
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

	return View{
		StructureID:      model.Item.ID,
		StructureShortID: model.Item.ShortID,
		StructureName:    model.Item.StructureName,
		ItemName:         model.Item.ItemName,
		ItemID:           model.Item.ItemID,
		ItemShortID:      model.Item.ItemShortID,
		ProjectID:        model.Item.ProjectID,
		Locale:           locale,
		Index:            model.Item.Index,
		Groups:           model.Item.Groups,
		Behaviour:        model.Item.Behaviour,
		Value:            model.Item.Value,
		Connections:      connections,
		CreatedAt:        model.Item.CreatedAt,
		UpdatedAt:        model.Item.UpdatedAt,
	}
}
