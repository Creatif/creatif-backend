package lists

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type DeleteRangeByID struct {
	Name      string   `json:"name"`
	ID        string   `json:"id"`
	ShortID   string   `json:"shortID"`
	Items     []string `json:"items"`
	ProjectID string   `param:"projectID"`
	Locale    string   `param:"locale"`
}

func SanitizeDeleteRangeByID(model DeleteRangeByID) DeleteRangeByID {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.ID = p.Sanitize(model.ID)
	model.ShortID = p.Sanitize(model.ShortID)

	model.Items = sdk.Map(model.Items, func(idx int, value string) string {
		return p.Sanitize(value)
	})

	return model
}
