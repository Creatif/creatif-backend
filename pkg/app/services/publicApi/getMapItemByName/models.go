package getMapItemByName

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Model struct {
	ProjectID   string
	Name        string
	Locale      string
	VersionName string
}

func NewModel(projectId, versionName, name, locale string) Model {
	return Model{
		ProjectID:   projectId,
		Name:        name,
		Locale:      locale,
		VersionName: versionName,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID":   a.ProjectID,
		"versionName": a.VersionName,
		"name":        a.Name,
		"locale":      a.Locale,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("versionName", validation.Required),
			validation.Key("locale", validation.RuneLength(3, 3), validation.By(func(value interface{}) error {
				l := value.(string)
				if l == "" {
					return nil
				}

				_, err := locales.GetIDWithAlpha(l)
				if err != nil {
					return errors.New(fmt.Sprintf("Invalid locale %s. This locale does not exist.", l))
				}
				return nil
			})),
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

	Name    string `json:"itemName"`
	ID      string `json:"itemId"`
	ShortID string `json:"itemShortId"`

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
	Item        Item
	Connections []ConnectionItem
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

	return View{
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
		CreatedAt:        model.Item.CreatedAt,
		UpdatedAt:        model.Item.UpdatedAt,
	}
}
