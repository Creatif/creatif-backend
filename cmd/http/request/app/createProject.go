package app

import "github.com/microcosm-cc/bluemonday"

type CreateProject struct {
	Name string `json:"name"`
}

func SanitizeProject(model CreateProject) CreateProject {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)

	return model
}
