package references

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type PaginateReferences struct {
	ProjectID         string `param:"projectID"`
	ParentID          string `param:"parentID"`
	ChildID           string `param:"childID"`
	StructureType     string `param:"structureType"`
	RelationshipType  string `param:"relationshipType"`
	ParentStructureID string `param:"parentStructureId"`
	ChildStructureID  string `param:"childStructureId"`
	Limit             int    `query:"limit"`
	Page              int    `query:"page"`
	Locales           string
	Search            string `query:"search"`
	OrderBy           string `query:"orderBy"`
	OrderDirection    string `query:"direction"`

	Filters   string `query:"filters"`
	Behaviour string `query:"behaviour"`
	Fields    string `query:"fields"`
	Groups    string `query:"groups"`
	In        string `query:"in"`

	SanitizedLocales []string
	SanitizedGroups  []string
	SanitizedFields  []string
}

func SanitizePaginateReferences(model PaginateReferences) PaginateReferences {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)

	model.ParentID = p.Sanitize(model.ParentID)
	model.ChildID = p.Sanitize(model.ChildID)
	model.StructureType = p.Sanitize(model.StructureType)
	model.RelationshipType = p.Sanitize(model.RelationshipType)
	model.ParentStructureID = p.Sanitize(model.ParentStructureID)
	model.ChildStructureID = p.Sanitize(model.ChildStructureID)

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

	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Search = p.Sanitize(model.Search)
	model.OrderDirection = p.Sanitize(model.OrderDirection)
	model.OrderBy = p.Sanitize(model.OrderBy)

	return model
}
