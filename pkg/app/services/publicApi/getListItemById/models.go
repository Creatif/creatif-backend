package getListItemById

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Model struct {
	ProjectID   string
	ItemID      string
	VersionName string

	Options Options
}

type Options struct {
	ValueOnly bool
}

func NewModel(versionName, projectId, itemId string, options Options) Model {
	return Model{
		ProjectID:   projectId,
		ItemID:      itemId,
		Options:     options,
		VersionName: versionName,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID":   a.ProjectID,
		"itemId":      a.ItemID,
		"versionName": a.VersionName,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("versionName", validation.When(a.VersionName != "", validation.RuneLength(1, 200))),
			validation.Key("projectID", validation.Required, validation.RuneLength(27, 27)),
			validation.Key("itemId", validation.Required, validation.RuneLength(27, 27)),
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

type View struct {
	StructureID      string `json:"structureId,omitempty"`
	StructureShortID string `json:"structureShortId,omitempty"`
	StructureName    string `json:"structureName,omitempty"`

	Name    string `json:"name,omitempty"`
	ID      string `json:"id,omitempty"`
	ShortID string `json:"shortId,omitempty"`

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
	Item    Item
	Options Options
}

func newView(model LogicModel) interface{} {
	if model.Options.ValueOnly {
		var m map[string]interface{}
		// ok to ignore
		json.Unmarshal(model.Item.Value, &m)

		return m
	}

	locale, _ := locales.GetAlphaWithID(model.Item.Locale)

	view := View{
		StructureID:      model.Item.StructureID,
		StructureShortID: model.Item.StructureShortID,
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
		CreatedAt:        nil,
		UpdatedAt:        nil,
	}

	if !model.Options.ValueOnly {
		view.CreatedAt = &model.Item.CreatedAt
		view.UpdatedAt = &model.Item.UpdatedAt
	}

	return view
}
