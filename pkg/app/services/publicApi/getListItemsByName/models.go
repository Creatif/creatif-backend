package getListItemsByName

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Model struct {
	ProjectID     string
	StructureName string
	Name          string
	Locale        string
	VersionName   string
	Options       Options
}

type Options struct {
	ValueOnly bool
}

func NewModel(versionName, projectId, structureName, name, locale string, options Options) Model {
	return Model{
		ProjectID:     projectId,
		StructureName: structureName,
		VersionName:   versionName,
		Name:          name,
		Locale:        locale,
		Options:       options,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID":     a.ProjectID,
		"name":          a.Name,
		"structureName": a.StructureName,
		"locale":        a.Locale,
		"versionName":   a.VersionName,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("versionName", validation.When(a.VersionName != "", validation.RuneLength(1, 200))),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("name", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("structureName", validation.Required, validation.RuneLength(1, 200)),
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

type ConnectionsView struct {
	Parents  []string `json:"parents"`
	Children []string `json:"children"`
}

type connections struct {
	parents  []string
	children []string
}

func newConnections() connections {
	return connections{
		parents:  []string{},
		children: []string{},
	}
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

	Connections ConnectionsView `json:"connections"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type LogicModel struct {
	Items       []Item
	Connections map[string]connections
	Options     Options
}

func newView(model LogicModel) interface{} {
	if model.Options.ValueOnly {
		returnValue := make([]map[string]interface{}, len(model.Items))
		for i, val := range model.Items {
			var m map[string]interface{}
			// ok to ignore
			json.Unmarshal(val.Value, &m)
			connections := model.Connections[val.ItemID]

			m["connections"] = ConnectionsView{
				Parents:  connections.parents,
				Children: connections.children,
			}

			returnValue[i] = m
		}

		return returnValue
	}

	views := make([]View, len(model.Items))
	for i, item := range model.Items {
		locale, _ := locales.GetAlphaWithID(item.Locale)
		connections := model.Connections[item.ItemID]

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
			Connections: ConnectionsView{
				Parents:  connections.parents,
				Children: connections.children,
			},
			CreatedAt: nil,
			UpdatedAt: nil,
		}

		if !model.Options.ValueOnly {
			views[i].CreatedAt = &item.CreatedAt
			views[i].UpdatedAt = &item.UpdatedAt
		}
	}

	return views
}
