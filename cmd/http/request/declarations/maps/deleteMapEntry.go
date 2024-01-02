package maps

import (
	"github.com/microcosm-cc/bluemonday"
)

type DeleteMapEntry struct {
	Name         string `param:"name"`
	VariableName string `param:"variableName"`
	ProjectID    string `param:"projectID"`
}

func SanitizeDeleteMapEntry(model DeleteMapEntry) DeleteMapEntry {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.VariableName = p.Sanitize(model.VariableName)

	return model
}
