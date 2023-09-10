package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetNode struct {
	ID     string   `param:"id"`
	Fields []string `param:"fields"`
}

func SanitizeGetNode(model GetNode) GetNode {
	p := bluemonday.StrictPolicy()
	model.ID = p.Sanitize(model.ID)
	if len(model.Fields) > 0 {
		newFields := make([]string, len(model.Fields))

		for i := 0; i < len(model.Fields); i++ {
			newFields[i] = p.Sanitize(model.Fields[i])
		}

		model.Fields = newFields
	}

	return model
}
