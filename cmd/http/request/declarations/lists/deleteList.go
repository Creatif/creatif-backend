package lists

import "github.com/microcosm-cc/bluemonday"

type DeleteList struct {
	Name      string `param:"name"`
	ID        string `json:"id"`
	ShortID   string `json:"shortID"`
	ProjectID string `param:"projectID"`
	Locale    string `param:"locale"`
}

func SanitizeDeleteList(model DeleteList) DeleteList {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ID = p.Sanitize(model.ID)
	model.ShortID = p.Sanitize(model.ShortID)
	model.Locale = p.Sanitize(model.Locale)

	return model
}
