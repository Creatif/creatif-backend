package publicApi

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetMapItemByNameOptions struct {
	ValueOnly bool
}

type GetMapItemByName struct {
	ProjectID     string `param:"projectId"`
	StructureName string `param:"structureName"`
	Name          string `param:"name"`
	Locale        string `query:"locale"`
	Options       string `query:"options"`

	ResolvedOptions GetListItemByIDOptions
}

func SanitizeGetMapItemByName(model GetMapItemByName) GetMapItemByName {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.StructureName = p.Sanitize(model.StructureName)
	model.Locale = p.Sanitize(model.Locale)
	model.Name = p.Sanitize(model.Name)
	model.Options = p.Sanitize(model.Options)

	if model.Options != "" {
		model.ResolvedOptions = resolveListOptions(model.Options)
	}

	return model
}
