package app

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type AddGroups struct {
	Groups    []string `json:"groups"`
	ProjectID string   `param:"projectId"`
}

func SanitizeAddGroups(model AddGroups) AddGroups {
	p := bluemonday.StrictPolicy()
	model.Groups = sdk.Map(model.Groups, func(idx int, value string) string {
		return p.Sanitize(value)
	})
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
