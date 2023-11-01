package lists

import "github.com/microcosm-cc/bluemonday"

type DeleteListItemByID struct {
	Name        string `param:"name"`
	ID          string `json:"id"`
	ShortID     string `json:"shortID"`
	ItemID      string `param:"itemID"`
	ItemShortID string `json:"itemShortID"`
	ProjectID   string `param:"projectID"`
	Locale      string `param:"locale"`
}

func SanitizeDeleteListItemByID(model DeleteListItemByID) DeleteListItemByID {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.ItemID = p.Sanitize(model.ItemID)
	model.ItemShortID = p.Sanitize(model.ItemShortID)
	model.ShortID = p.Sanitize(model.ShortID)

	return model
}
