package maps

import (
	"github.com/microcosm-cc/bluemonday"
)

type QueryMapVariable struct {
	Name      string `param:"name"`
	ItemID    string `param:"itemID"`
	ProjectID string `param:"projectID"`
}

func SanitizeQueryMapVariable(model QueryMapVariable) QueryMapVariable {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ItemID = p.Sanitize(model.ItemID)

	return model
}
