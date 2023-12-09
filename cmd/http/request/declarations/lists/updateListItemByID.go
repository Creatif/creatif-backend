package lists

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type UpdateListItemByIDValues struct {
	Name      string   `json:"name"`
	Metadata  string   `json:"metadata"`
	Locale    string   `json:"locale"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Value     string   `json:"value"`
}

type UpdateListItemByID struct {
	Name      string                   `param:"name"`
	ItemID    string                   `param:"itemID"`
	Values    UpdateListItemByIDValues `json:"values"`
	ProjectID string                   `param:"projectID"`
	Fields    string                   `query:"fields"`

	ResolvedFields []string
}

func SanitizeUpdateListItemByID(model UpdateListItemByID) UpdateListItemByID {
	p := bluemonday.StrictPolicy()

	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ItemID = p.Sanitize(model.ItemID)
	model.Fields = p.Sanitize(model.Fields)

	model.ResolvedFields = sdk.Map(strings.Split(model.Fields, "|"), func(idx int, value string) string {
		trimmed := strings.Trim(value, " ")
		return trimmed
	})

	model.Values = UpdateListItemByIDValues{
		Name:      p.Sanitize(model.Values.Name),
		Locale:    p.Sanitize(model.Values.Locale),
		Behaviour: p.Sanitize(model.Values.Behaviour),
		Groups: sdk.Map(model.Values.Groups, func(idx int, value string) string {
			return p.Sanitize(value)
		}),
		Metadata: model.Values.Metadata,
		Value:    model.Values.Value,
	}

	return model
}
