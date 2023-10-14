package lists

import "github.com/microcosm-cc/bluemonday"

type SwitchByIndex struct {
	Name        string `param:"name"`
	Source      int64  `param:"source"`
	Destination int64  `param:"destination"`
	ProjectID   string `param:"projectID"`
	Locale      string `param:"locale"`
}

func SanitizeSwitchByIndex(model SwitchByIndex) SwitchByIndex {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)

	return model
}
