package lists

import (
	"github.com/microcosm-cc/bluemonday"
)

type QueryListByID struct {
	Name               string `param:"name"`
	ItemID             string `param:"itemID"`
	ProjectID          string `param:"projectID"`
	ConnectionViewType string `query:"connectionViewType"`
}

func SanitizeQueryListByID(model QueryListByID) QueryListByID {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ItemID = p.Sanitize(model.ItemID)
	model.ConnectionViewType = p.Sanitize(model.ConnectionViewType)

	return model
}
