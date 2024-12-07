package connections

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type PaginateConnections struct {
	ProjectID        string `param:"projectID"`
	StructureType    string `param:"structureType"`
	StructureID      string `param:"structureId"`
	ParentVariableID string `param:"parentVariableId"`
	Locales          string `query:"locales"`
	Limit            int    `query:"limit"`
	Page             int    `query:"page"`
	Filters          string `query:"filters"`
	Behaviour        string `query:"behaviour"`
	Fields           string `query:"fields"`
	Groups           string `query:"groups"`
	Search           string `query:"search"`
	OrderBy          string `query:"orderBy"`
	OrderDirection   string `query:"direction"`
	In               string `query:"in"`

	SanitizedGroups  []string
	SanitizedLocales []string
	SanitizedFields  []string
}

func SanitizePaginateConnections(model PaginateConnections) PaginateConnections {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.StructureType = p.Sanitize(model.StructureType)
	model.StructureID = p.Sanitize(model.StructureID)
	model.ParentVariableID = p.Sanitize(model.ParentVariableID)
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Search = p.Sanitize(model.Search)
	model.OrderDirection = p.Sanitize(model.OrderDirection)
	model.In = p.Sanitize(model.In)
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Behaviour = p.Sanitize(model.Behaviour)

	if model.Groups != "" {
		model.SanitizedGroups = sdk.Map(strings.Split(model.Groups, ","), func(idx int, value string) string {
			return p.Sanitize(strings.TrimSpace(value))
		})
	}

	if model.Fields != "" {
		model.SanitizedFields = sdk.Map(strings.Split(model.Fields, ","), func(idx int, value string) string {
			return p.Sanitize(strings.TrimSpace(value))
		})
	}

	if model.Locales != "" {
		model.SanitizedLocales = sdk.Map(strings.Split(model.Locales, ","), func(idx int, value string) string {
			sanitized := p.Sanitize(strings.TrimSpace(value))
			locale, _ := locales.GetIDWithAlpha(sanitized)

			return locale
		})
	}

	return model
}
