package lists

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type UpdateListValues struct {
	Name string `json:"name"`
}

type UpdateList struct {
	Fields    []string
	Name      string
	Values    UpdateListValues
	ProjectID string
	Locale    string
}

func SanitizeUpdateList(model UpdateList) UpdateList {
	p := bluemonday.StrictPolicy()

	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.Fields = sdk.Map(model.Fields, func(idx int, value string) string {
		return p.Sanitize(value)
	})
	model.Values = UpdateListValues{Name: p.Sanitize(model.Values.Name)}

	return model
}
