package publicApi

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetMapItemByIDOptions struct {
	ValueOnly bool
}

type GetMapItemByID struct {
	ProjectID string `param:"projectId"`
	ItemID    string `param:"id"`
	Options   string `query:"options"`

	ResolvedOptions GetMapItemByIDOptions
}

func SanitizeGetMapItemByID(model GetMapItemByID) GetMapItemByID {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ItemID = p.Sanitize(model.ItemID)

	if model.Options != "" {
		model.ResolvedOptions = resolveMapOptions(model.Options)
	}

	return model
}
