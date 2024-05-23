package publicApi

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetListItemByIDOptions struct {
	ValueOnly bool
}

type GetListItemByID struct {
	ProjectID   string `param:"projectId"`
	ItemID      string `param:"id"`
	VersionName string
	Options     string `query:"options"`

	ResolvedOptions GetListItemByIDOptions
}

func SanitizeGetListItemByID(model GetListItemByID) GetListItemByID {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ItemID = p.Sanitize(model.ItemID)
	model.Options = p.Sanitize(model.Options)
	model.VersionName = p.Sanitize(model.VersionName)

	if model.Options != "" {
		model.ResolvedOptions = resolveListOptions(model.Options)
	}

	return model
}
