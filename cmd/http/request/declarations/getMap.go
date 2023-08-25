package declarations

import (
	"github.com/microcosm-cc/bluemonday"
)

type GetMap struct {
	ID     string `param:"id"`
	Return string `query:"return"`
	Fields string `query:"fields"`
}

func SanitizeGetMap(model GetMap) GetMap {
	p := bluemonday.StrictPolicy()
	model.ID = p.Sanitize(model.ID)
	model.Return = p.Sanitize(model.Return)
	model.Fields = p.Sanitize(model.Fields)

	return model
}
