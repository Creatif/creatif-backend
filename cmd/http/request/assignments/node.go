package assignments

import "github.com/microcosm-cc/bluemonday"

type AssignNodeTextValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func SanitizeTextValue(model AssignNodeTextValue) AssignNodeTextValue {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)

	return model
}
