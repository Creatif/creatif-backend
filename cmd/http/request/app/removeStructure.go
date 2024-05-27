package app

import "github.com/microcosm-cc/bluemonday"

type RemoveStructure struct {
	ProjectID string `param:"projectId"`
	ID        string `json:"id"`
	Type      string `json:"type"`
}

func SanitizeRemoveStructure(model RemoveStructure) RemoveStructure {
	p := bluemonday.StrictPolicy()
	model.ID = p.Sanitize(model.ID)
	model.Type = p.Sanitize(model.Type)
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
