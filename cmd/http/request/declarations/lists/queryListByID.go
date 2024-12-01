package lists

import (
	"github.com/microcosm-cc/bluemonday"
)

type QueryListByID struct {
	Name                    string `param:"name"`
	ItemID                  string `param:"itemID"`
	ProjectID               string `param:"projectID"`
	ConnectionReplaceMethod string `query:"connectionReplaceMethod"`
}

func SanitizeQueryListByID(model QueryListByID) QueryListByID {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ItemID = p.Sanitize(model.ItemID)
	model.ConnectionReplaceMethod = p.Sanitize(model.ConnectionReplaceMethod)

	return model
}
