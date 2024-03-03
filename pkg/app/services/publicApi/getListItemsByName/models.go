package getListItemsByName

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
	Items       []Item
	Connections map[string][]ConnectionItem
}

func newView(model LogicModel) []View {
	views := make([]View, len(model.Items))
	for i, item := range model.Items {
		locale, _ := locales.GetAlphaWithID(item.Locale)
		connectionViews := make(map[string]ConnectionView)

		connections, ok := model.Connections[item.ItemID]
		if ok {
			for _, c := range connections {
				connectionLocale, _ := locales.GetAlphaWithID(item.Locale)

				connectionViews[c.ConnectionName] = ConnectionView{
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
		}

		views[i] = View{
			StructureID:      item.ID,
			StructureShortID: item.ShortID,
			StructureName:    item.StructureName,
			Name:             item.ItemName,
			ID:               item.ItemID,
			ShortID:          item.ItemShortID,
			ProjectID:        item.ProjectID,
			Locale:           locale,
			Index:            item.Index,
			Groups:           item.Groups,
			Behaviour:        item.Behaviour,
			Value:            item.Value,
			Connections:      connectionViews,
			CreatedAt:        item.CreatedAt,
			UpdatedAt:        item.UpdatedAt,
		}
	}

	return views
}