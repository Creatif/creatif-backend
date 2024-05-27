package app

import "github.com/microcosm-cc/bluemonday"

type TruncateStructure struct {
	ProjectID string `param:"projectId"`
	ID        string `json:"id"`
	Type      string `json:"type"`
}

func SanitizeTruncateStructure(model TruncateStructure) TruncateStructure {
	p := bluemonday.StrictPolicy()
	model.ID = p.Sanitize(model.ID)
	model.Type = p.Sanitize(model.Type)
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
