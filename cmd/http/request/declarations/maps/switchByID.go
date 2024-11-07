package maps

import "github.com/microcosm-cc/bluemonday"

type SwitchByID struct {
	Name           string `param:"name"`
	Source         string `param:"source"`
	Destination    string `param:"destination"`
	ProjectID      string `param:"projectID"`
	OrderDirection string `param:"orderDirection"`
}

func SanitizeSwitchByID(model SwitchByID) SwitchByID {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.Source = p.Sanitize(model.Source)
	model.Destination = p.Sanitize(model.Destination)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.OrderDirection = p.Sanitize(model.OrderDirection)

	return model
}
