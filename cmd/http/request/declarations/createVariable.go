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
}

type ValidationLength struct {
	Min   int `json:"min"`
	Max   int `json:"max"`
	Exact int `json:"exact"`
}

func SanitizeVariable(model CreateVariable) CreateVariable {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.Behaviour = p.Sanitize(model.Behaviour)
	model.Metadata = p.Sanitize(model.Metadata)
	model.Value = p.Sanitize(model.Value)

	model.Groups = sdk.Sanitize(model.Groups, func(k int, v string) string {
		return p.Sanitize(v)
	})

	return model
}
