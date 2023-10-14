package lists

import "github.com/microcosm-cc/bluemonday"

type DeleteList struct {
	Name      string
	ProjectID string
	Locale    string
}

func SanitizeDeleteList(model DeleteList) DeleteList {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)

	return model
}
