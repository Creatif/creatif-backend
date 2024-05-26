package getManyItems

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Model struct {
	ProjectID   string
	VersionName string
	IDs         []string
	Options     Options
}

type Options struct {
	ValueOnly bool
}

func NewModel(versionName, projectId string, ids []string, options Options) Model {
	return Model{
		ProjectID:   projectId,
		VersionName: versionName,
		IDs:         ids,
		Options:     options,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID":   a.ProjectID,
		"versionName": a.VersionName,
		"ids":         nil,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("versionName", validation.When(a.VersionName != "", validation.RuneLength(1, 200))),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("ids", validation.By(func(value interface{}) error {
				if len(a.IDs) > 100 {
					return errors.New("Maximum number of ids is 100")
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
