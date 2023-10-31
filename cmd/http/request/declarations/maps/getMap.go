package maps

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type GetMap struct {
	Name      string `param:"name"`
	ID        string `json:"id"`
	ShortID   string `json:"shortID"`
	Fields    string `query:"fields"`
	ProjectID string `param:"projectID"`
	Locale    string `param:"locale"`
	Groups    string `query:"groups"`

	SanitizedGroups []string
}

func SanitizeGetMap(model GetMap) GetMap {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Fields = p.Sanitize(model.Fields)
	model.Locale = p.Sanitize(model.Locale)
	model.ShortID = p.Sanitize(model.ShortID)
	model.ID = p.Sanitize(model.ID)

	if model.Groups != "" {
		newGroups := sdk.Map(strings.Split(model.Groups, ","), func(idx int, value string) string {
			return p.Sanitize(strings.TrimSpace(value))
		})

		model.SanitizedGroups = newGroups
	}

	return model
}
