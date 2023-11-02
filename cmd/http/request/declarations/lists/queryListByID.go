package lists

import (
	"github.com/microcosm-cc/bluemonday"
)

type QueryListByID struct {
	Name        string `param:"name"`
	ShortID     string `json:"shortID"`
	ID          string `json:"id"`
	ItemID      string `json:"itemID"`
	ItemShortID string `json:"itemShortID"`
	ProjectID   string `param:"projectID"`
	Locale      string `param:"locale"`
}

func SanitizeQueryListByID(model QueryListByID) QueryListByID {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.Locale = p.Sanitize(model.Locale)
	model.ID = p.Sanitize(model.ID)
	model.ShortID = p.Sanitize(model.ShortID)
	model.ItemID = p.Sanitize(model.ItemID)
	model.ItemShortID = p.Sanitize(model.ItemShortID)

	return model
}
