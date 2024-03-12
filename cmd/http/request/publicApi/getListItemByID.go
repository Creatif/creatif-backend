package publicApi

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetListItemByID struct {
	ProjectID string `param:"projectId"`
	ItemID    string `param:"id"`
}

func SanitizeGetListItemByID(model GetListItemByID) GetListItemByID {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ItemID = p.Sanitize(model.ItemID)

	return model
}
