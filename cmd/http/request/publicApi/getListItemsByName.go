package publicApi

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetListItemByNameOptions struct {
	ValueOnly bool
}

type GetListItemByName struct {
	ProjectID     string `param:"projectId"`
	StructureName string `param:"structureName"`
	Name          string `param:"name"`
	Locale        string `query:"locale"`
	VersionName   string
	Options       string `query:"options"`

	ResolvedOptions GetListItemByIDOptions
}

func SanitizeGetListItemByName(model GetListItemByName) GetListItemByName {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.StructureName = p.Sanitize(model.StructureName)
	model.Locale = p.Sanitize(model.Locale)
	model.Name = p.Sanitize(model.Name)
	model.Options = p.Sanitize(model.Options)
	model.VersionName = p.Sanitize(model.VersionName)

	if model.Options != "" {
		model.ResolvedOptions = resolveListOptions(model.Options)
	}

	return model
}
