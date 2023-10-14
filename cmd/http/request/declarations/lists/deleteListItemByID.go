package lists

import "github.com/microcosm-cc/bluemonday"

type DeleteListItemByID struct {
	Name      string `param:"name"`
	ItemID    string `param:"itemID"`
	ProjectID string `param:"projectID"`
	Locale    string `param:"locale"`
}

func SanitizeDeleteListItem(model DeleteListItemByID) DeleteListItemByID {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.ItemID = p.Sanitize(model.ItemID)

	return model
}
