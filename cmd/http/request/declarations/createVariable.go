package declarations

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type CreateVariable struct {
	Name      string   `json:"name"`
	Groups    []string `json:"groups"`
	Behaviour string   `json:"behaviour"`
	Metadata  string   `json:"metadata"`
	Value     string   `json:"value"`
	ProjectID string   `json:"projectID"`
	Language  string   `json:"language"`
}

func SanitizeVariable(model CreateVariable) CreateVariable {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.Behaviour = p.Sanitize(model.Behaviour)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Language = p.Sanitize(model.Language)

	model.Groups = sdk.Sanitize(model.Groups, func(k int, v string) string {
		return p.Sanitize(v)
	})

	return model
}
