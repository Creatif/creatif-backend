package lists

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type UpdateListItemByIDValues struct {
	Name      string   `json:"name"`
	Metadata  string   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Value     string   `json:"value"`
}

type UpdateListItemByID struct {
	Fields      []string                 `query:"projectID"`
	ListName    string                   `param:"name"`
	ListID      string                   `json:"id"`
	ListShortID string                   `json:"shortID"`
	Locale      string                   `param:"locale"`
	ItemID      string                   `param:"itemID"`
	ItemShortID string                   `param:"itemShortID"`
	Values      UpdateListItemByIDValues `json:"values"`
	ProjectID   string                   `param:"projectID"`
}

func SanitizeUpdateListItemByID(model UpdateListItemByID) UpdateListItemByID {
	p := bluemonday.StrictPolicy()

	model.ListName = p.Sanitize(model.ListName)
	model.ListID = p.Sanitize(model.ListID)
	model.ListShortID = p.Sanitize(model.ListShortID)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.ItemID = p.Sanitize(model.ItemID)
	model.ItemShortID = p.Sanitize(model.ItemShortID)
	model.Fields = sdk.Map(model.Fields, func(idx int, value string) string {
		return p.Sanitize(value)
	})
	model.Values = UpdateListItemByIDValues{
		Name:      p.Sanitize(model.Values.Name),
		Behaviour: p.Sanitize(model.Values.Behaviour),
		Groups: sdk.Map(model.Values.Groups, func(idx int, value string) string {
			return p.Sanitize(value)
		}),
		Metadata: model.Values.Metadata,
		Value:    model.Values.Value,
	}

	return model
}
