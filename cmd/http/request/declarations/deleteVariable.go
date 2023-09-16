package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type DeleteVariable struct {
	Name string `json:"name"`
}

func SanitizeDeleteVariable(model DeleteVariable) DeleteVariable {
	p := bluemonday.StrictPolicy()
	model.Name = p.Sanitize(model.Name)

	return model
}
