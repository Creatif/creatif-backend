package maps

import (
	"github.com/microcosm-cc/bluemonday"
)

type QueryMapVariable struct {
	Name                    string `param:"name"`
	ItemID                  string `param:"itemID"`
	ProjectID               string `param:"projectID"`
	ConnectionReplaceMethod string `query:"connectionReplaceMethod"`
}

func SanitizeQueryMapVariable(model QueryMapVariable) QueryMapVariable {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ItemID = p.Sanitize(model.ItemID)
	model.ConnectionReplaceMethod = p.Sanitize(model.ConnectionReplaceMethod)

	return model
}
