package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetBatchedVariables struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func SanitizeGetBatchedVariables(model []GetBatchedVariables) []GetBatchedVariables {
	p := bluemonday.StrictPolicy()
	sanitized := make([]GetBatchedVariables, 0)

	for _, n := range model {
		sanitized = append(sanitized, GetBatchedVariables{
			Name: p.Sanitize(n.Name),
			Type: p.Sanitize(n.Type),
		})
	}

	return sanitized
}
