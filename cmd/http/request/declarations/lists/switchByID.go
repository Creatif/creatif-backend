package lists

import "github.com/microcosm-cc/bluemonday"

type SwitchByID struct {
	Name        string `param:"name"`
	ID          string `json:"id"`
	ShortID     string `json:"shortID"`
	Source      string `param:"source"`
	Destination string `param:"destination"`
	ProjectID   string `param:"projectID"`
}

func SanitizeSwitchByID(model SwitchByID) SwitchByID {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.Source = p.Sanitize(model.Source)
	model.Destination = p.Sanitize(model.Destination)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ID = p.Sanitize(model.ID)
	model.ShortID = p.Sanitize(model.ShortID)

	return model
}
