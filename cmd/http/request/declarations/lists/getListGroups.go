package lists

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetListGroups struct {
	Name      string `param:"name"`
	ShortID   string `json:"shortID"`
	ID        string `json:"id"`
	ProjectID string `param:"projectID"`
	Locale    string `param:"locale"`
}

func SanitizeGetListGroups(model GetListGroups) GetListGroups {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ID = p.Sanitize(model.ID)
	model.ShortID = p.Sanitize(model.ShortID)
	model.Locale = p.Sanitize(model.Locale)

	return model
}
