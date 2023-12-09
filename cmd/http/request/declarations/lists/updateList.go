package lists

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type UpdateListValues struct {
	Name string `json:"name"`
}

type UpdateList struct {
	Fields    []string         `json:"fields"`
	Name      string           `param:"name"`
	ID        string           `json:"id"`
	ShortID   string           `json:"shortID"`
	Values    UpdateListValues `json:"values"`
	ProjectID string           `param:"projectID"`
}

func SanitizeUpdateList(model UpdateList) UpdateList {
	p := bluemonday.StrictPolicy()

	model.Name = p.Sanitize(model.Name)
	model.ID = p.Sanitize(model.ID)
	model.ShortID = p.Sanitize(model.ShortID)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Fields = sdk.Map(model.Fields, func(idx int, value string) string {
		return p.Sanitize(value)
	})
	model.Values = UpdateListValues{Name: p.Sanitize(model.Values.Name)}

	return model
}
