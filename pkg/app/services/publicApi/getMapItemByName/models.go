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
	ProjectID     string
	Name          string
	Locale        string
	StructureName string
	Options       Options
	VersionName   string
}

type Options struct {
	ValueOnly bool
}

func NewModel(versionName, projectId, structureName, name, locale string, options Options) Model {
	return Model{
		ProjectID:     projectId,
		Name:          name,
		StructureName: structureName,
		Locale:        locale,
		Options:       options,
		VersionName:   versionName,
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

	Connections ConnectionsView `json:"connections"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LogicModel struct {
	Item        Item
	Connections connections
	Options     Options
}

func newView(model LogicModel) interface{} {
	if model.Options.ValueOnly {
		return model.Item.Value
	}

	locale, _ := locales.GetAlphaWithID(model.Item.Locale)

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
		Connections: ConnectionsView{
			Parents:  model.Connections.parents,
			Children: model.Connections.children,
		},
		CreatedAt: model.Item.CreatedAt,
		UpdatedAt: model.Item.UpdatedAt,
	}
}
