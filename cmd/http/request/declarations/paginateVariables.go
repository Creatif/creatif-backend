package declarations

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

type PaginateVariables struct {
	NextID    string `query:"nextId"`
	PrevID    string `query:"prevId"`
	Field     string `query:"field"`
	OrderBy   string `query:"orderBy"`
	Direction string `query:"direction"`
	Groups    string `query:"groups"`
	Limit     int    `query:"limit"`

	SanitizedGroups []string
}

func SanitizePaginateVariables(model PaginateVariables) PaginateVariables {
	p := bluemonday.StrictPolicy()
	model.NextID = p.Sanitize(model.NextID)
	model.PrevID = p.Sanitize(model.PrevID)
	model.Field = p.Sanitize(model.Field)
	model.OrderBy = p.Sanitize(model.OrderBy)
	model.Direction = p.Sanitize(model.Direction)
	model.Groups = p.Sanitize(model.Groups)

	if model.Groups != "" {
		newGroups := sdk.Map(strings.Split(model.Groups, ","), func(idx int, value string) string {
			return p.Sanitize(strings.TrimSpace(value))
		})

		model.SanitizedGroups = newGroups
	}

	return model
}
