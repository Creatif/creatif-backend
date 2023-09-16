package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetValue struct {
	Name string `json:"name"`
}

func SanitizeGetValue(model GetValue) GetValue {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)

	return model
}
