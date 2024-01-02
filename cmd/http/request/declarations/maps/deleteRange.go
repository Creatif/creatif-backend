package maps

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type DeleteRange struct {
	Name      string   `param:"name"`
	Items     []string `json:"items"`
	ProjectID string   `param:"projectID"`
}

func SanitizeDeleteRange(model DeleteRange) DeleteRange {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)

	model.Items = sdk.Map(model.Items, func(idx int, value string) string {
		return p.Sanitize(value)
	})

	return model
}
