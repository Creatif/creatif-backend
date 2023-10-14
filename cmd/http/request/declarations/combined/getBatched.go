package combined

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetBatchedStructures struct {
	ProjectID  string              `param:"projectID"`
	Locale     string              `param:"locale"`
	Structures []BatchedStructures `json:"structures"`
}

type BatchedStructures struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func SanitizeGetBatchedVariables(model GetBatchedStructures) GetBatchedStructures {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)

	sanitized := make([]BatchedStructures, 0)
	for _, n := range model.Structures {
		sanitized = append(sanitized, BatchedStructures{
			Name: p.Sanitize(n.Name),
			Type: p.Sanitize(n.Type),
		})
	}

	model.Structures = sanitized

	return model
}
