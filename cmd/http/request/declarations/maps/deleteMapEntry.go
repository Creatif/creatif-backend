package maps

import (
	"github.com/microcosm-cc/bluemonday"
)

type DeleteMapEntry struct {
	Name       string `param:"name"`
	MapID      string `json:"mapID"`
	MapShortID string `json:"mapShortID"`

	VariableName    string `param:"variableName"`
	VariableID      string `param:"variableID"`
	VariableShortID string `param:"variableShortID"`
	ProjectID       string `param:"projectID"`
	Locale          string `param:"locale"`
}

func SanitizeDeleteMapEntry(model DeleteMapEntry) DeleteMapEntry {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.VariableName = p.Sanitize(model.VariableName)
	model.VariableID = p.Sanitize(model.VariableID)
	model.VariableShortID = p.Sanitize(model.VariableShortID)
	model.Locale = p.Sanitize(model.Locale)
	model.MapID = p.Sanitize(model.MapID)
	model.MapShortID = p.Sanitize(model.MapShortID)

	return model
}
