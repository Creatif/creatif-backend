package lists

import "github.com/microcosm-cc/bluemonday"

type SwitchByID struct {
	// this can be project name
	Name        string
	Source      string
	Destination string
	ProjectID   string
	Locale      string
}

func SanitizeSwitchByID(model SwitchByID) SwitchByID {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.Source = p.Sanitize(model.Source)
	model.Destination = p.Sanitize(model.Destination)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)

	return model
}
