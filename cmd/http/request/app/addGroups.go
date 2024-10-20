package app

import (
	"creatif/pkg/lib/sdk"
	"github.com/microcosm-cc/bluemonday"
)

type SingleGroup struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Action string `json:"action"`
}

type AddGroups struct {
	Groups    []SingleGroup `json:"groups"`
	ProjectID string        `param:"projectId"`
}

func SanitizeAddGroups(model AddGroups) AddGroups {
	p := bluemonday.StrictPolicy()
	model.Groups = sdk.Map(model.Groups, func(idx int, value SingleGroup) SingleGroup {
		return SingleGroup{
			ID:     value.ID,
			Name:   value.Name,
			Type:   value.Type,
			Action: value.Action,
		}
	})
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
