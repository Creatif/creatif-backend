package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type DeleteMapEntry struct {
	Name      string `param:"name"`
	EntryName string `param:"entryName"`
	ProjectID string `param:"projectID"`
	Locale    string `json:"locale"`
}

func SanitizeDeleteMapEntry(model DeleteMapEntry) DeleteMapEntry {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.EntryName = p.Sanitize(model.EntryName)
	model.Locale = p.Sanitize(model.Locale)

	return model
}
