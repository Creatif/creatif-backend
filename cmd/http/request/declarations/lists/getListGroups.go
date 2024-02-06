package lists

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetListGroups struct {
	Name      string `param:"name"`
	ProjectID string `param:"projectID"`
	ItemID    string `param:"itemId"`
}

func SanitizeGetListGroups(model GetListGroups) GetListGroups {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ItemID = p.Sanitize(model.ItemID)

	return model
}
