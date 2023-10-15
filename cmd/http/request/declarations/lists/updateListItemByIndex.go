package lists

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type UpdateListItemByIndexValues struct {
	Name      string   `json:"name"`
	Metadata  string   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Value     string   `json:"value"`
}

type UpdateListItemByIndex struct {
	Fields    []string                    `query:"projectID"`
	ListName  string                      `param:"listName"`
	Locale    string                      `param:"locale"`
	Index     string                      `param:"index"`
	Values    UpdateListItemByIndexValues `json:"values"`
	ProjectID string                      `param:"projectID"`
}

func SanitizeUpdateListItemByIndex(model UpdateListItemByIndex) UpdateListItemByIndex {
	p := bluemonday.StrictPolicy()

	model.ListName = p.Sanitize(model.ListName)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Index = p.Sanitize(model.Index)
	model.Locale = p.Sanitize(model.Locale)
	model.Fields = sdk.Map(model.Fields, func(idx int, value string) string {
		return p.Sanitize(value)
	})
	model.Values = UpdateListItemByIndexValues{
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
