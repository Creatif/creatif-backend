package paginateListItems

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
	"time"
)

type Model struct {
	ProjectID     string
	StructureName string
	VersionName   string
	Options       Options

	Page    int
	Order   string
	SortBy  string
	Search  string
	Locales []string
	Groups  []string
}

type Options struct {
	ValueOnly bool
}

func NewModel(versionName, projectId, structureName string, page int, order string, sortBy, search string, lcls, groups []string, options Options) Model {
	return Model{
		StructureName: structureName,
		Options:       options,
		VersionName:   versionName,
		ProjectID:     projectId,
		Page:          page,
		Order:         order,
		SortBy:        sortBy,
		Search:        search,
		Locales:       lcls,
		Groups:        groups,
	}
}

func (a Model) Validate() map[string]string {
	v := map[string]interface{}{
		"projectID":     a.ProjectID,
		"order":         a.Order,
		"sortBy":        a.SortBy,
		"structureName": a.StructureName,
		"locales":       a.Locales,
		"versionName":   a.VersionName,
	}

	if err := validation.Validate(v,
		validation.Map(
			validation.Key("versionName", validation.When(a.VersionName != "", validation.RuneLength(1, 200))),
			validation.Key("structureName", validation.Required, validation.RuneLength(1, 200)),
			validation.Key("projectID", validation.Required, validation.RuneLength(26, 26)),
			validation.Key("order", validation.By(func(value interface{}) error {
				order := strings.ToLower(a.Order)
				if order != "desc" && order != "asc" {
					return errors.New("Order must be either DESC or ASC")
				}
				return nil
			})),
			validation.Key("sortBy", validation.By(func(value interface{}) error {
				sortBy := strings.ToLower(a.SortBy)
				if sortBy != "name" && sortBy != "created_at" && sortBy != "updated_at" && sortBy != "index" {
					return errors.New("Invalid sortBy field. sortBy can be: name, created_at, updated_at, index")
				}
				return nil
			})),
			validation.Key("locales", validation.By(func(value interface{}) error {
				for _, l := range a.Locales {
					_, err := locales.GetIDWithAlpha(l)
					if err != nil {
						return errors.New(fmt.Sprintf("Locale %s does not exist.", l))
					}
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
	Options     Options
}

func newView(model LogicModel) interface{} {
	if model.Options.ValueOnly {
		items := make([]interface{}, len(model.Items))

		for i, item := range model.Items {
			items[i] = item.Value
		}

		return items
	}

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
