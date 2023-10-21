package lists

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type ReplaceListVariable struct {
	Name      string   `json:"name"`
	Metadata  string   `json:"metadata"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Value     string   `json:"value"`
}

type ReplaceListItem struct {
	Name      string              `param:"name"`
	ItemName  string              `param:"itemName"`
	ProjectID string              `param:"projectID"`
	Locale    string              `param:"locale"`
	Variable  ReplaceListVariable `json:"variable"`
}

func SanitizeReplaceListItem(model ReplaceListItem) ReplaceListItem {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ItemName = p.Sanitize(model.ItemName)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)

	model.Variable = ReplaceListVariable{
		Name:     p.Sanitize(model.Variable.Name),
		Metadata: model.Variable.Metadata,
		Groups: sdk.Map(model.Variable.Groups, func(idx int, value string) string {
			return p.Sanitize(value)
		}),
		Behaviour: p.Sanitize(model.Variable.Behaviour),
		Value:     model.Variable.Value,
	}

	return model
}
