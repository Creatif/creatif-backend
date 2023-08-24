package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetNode struct {
	ID string `param:"id"`
}

func SanitizeGetNode(model GetNode) GetNode {
	p := bluemonday.StrictPolicy()
	model.ID = p.Sanitize(model.ID)

	return model
}
