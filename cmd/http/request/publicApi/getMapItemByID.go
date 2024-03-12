package publicApi

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetMapItemByID struct {
	ProjectID string `param:"projectId"`
	ItemID    string `param:"id"`
}

func SanitizeGetMapItemByID(model GetMapItemByID) GetMapItemByID {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ItemID = p.Sanitize(model.ItemID)

	return model
}
