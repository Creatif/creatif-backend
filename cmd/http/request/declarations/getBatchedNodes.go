package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetBatchedNodes struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func SanitizeGetBatchedNodes(model []GetBatchedNodes) []GetBatchedNodes {
	p := bluemonday.StrictPolicy()
	sanitized := make([]GetBatchedNodes, 0)

	for _, n := range model {
		sanitized = append(sanitized, GetBatchedNodes{
			Name: p.Sanitize(n.Name),
			Type: p.Sanitize(n.Type),
		})
	}

	return sanitized
}
