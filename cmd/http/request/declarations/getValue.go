package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetValue struct {
	Name      string `param:"name"`
	ProjectID string `param:"projectID"`
}

func SanitizeGetValue(model GetValue) GetValue {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)

	return model
}
