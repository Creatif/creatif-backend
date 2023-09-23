package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetBatchedStructures struct {
	ProjectID  string              `param:"projectID"`
	Structures []BatchedStructures `json:"structures"`
}

type BatchedStructures struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func SanitizeGetBatchedVariables(model GetBatchedStructures) GetBatchedStructures {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)

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
