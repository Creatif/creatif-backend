package lists

import "github.com/microcosm-cc/bluemonday"

type DeleteListItemByIndex struct {
	Name      string `param:"name"`
	ItemIndex string `param:"itemIndex"`
	ProjectID string `param:"projectID"`
	Locale    string `param:"locale"`
}

func SanitizeDeleteListItemByIndex(model DeleteListItemByIndex) DeleteListItemByIndex {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.ItemIndex = p.Sanitize(model.ItemIndex)

	return model
}
