package toggleProduction

import (
	"github.com/microcosm-cc/bluemonday"
)

type ToggleProduction struct {
	ProjectID string `param:"projectId"`
	ID        string `param:"id"`
}

func SanitizeToggleProduction(model ToggleProduction) ToggleProduction {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)
	model.ID = p.Sanitize(model.ID)

	return model
}
