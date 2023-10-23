package lists

import "github.com/microcosm-cc/bluemonday"

type DeleteListItemByIndex struct {
	Name      string `param:"name"`
	ItemIndex int64  `param:"itemIndex"`
	ProjectID string `param:"projectID"`
	Locale    string `param:"locale"`
}

func SanitizeDeleteListItemByIndex(model DeleteListItemByIndex) DeleteListItemByIndex {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)

	return model
}
