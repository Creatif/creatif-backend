package publicApi

import (
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type GetManyOptions struct {
	ValueOnly bool
}

type GetMany struct {
	ProjectID   string `param:"projectId"`
	VersionName string
	Options     string `query:"options"`
	IDs         string `query:"ids"`

	ResolvedOptions GetListItemByIDOptions
	ResolvedIds     []string
}

func SanitizeGetMany(model GetMany) GetMany {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.IDs = p.Sanitize(model.IDs)
	model.Options = p.Sanitize(model.Options)
	model.VersionName = p.Sanitize(model.VersionName)

	if model.Options != "" {
		model.ResolvedOptions = resolveListOptions(model.Options)
	}

	model.ResolvedIds = strings.Split(model.IDs, ",")

	return model
}
