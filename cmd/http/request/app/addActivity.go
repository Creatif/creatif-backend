package app

import "github.com/microcosm-cc/bluemonday"

type AddActivity struct {
	ProjectID string `param:"projectId"`
	Data      string `json:"data"`
}

func SanitizeAddActivity(model AddActivity) AddActivity {
	p := bluemonday.StrictPolicy()
	model.ProjectID = p.Sanitize(model.ProjectID)

	return model
}
