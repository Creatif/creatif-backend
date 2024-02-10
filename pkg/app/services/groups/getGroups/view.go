package getGroups

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
)

func newView(model []declarations.Group) []string {
	return sdk.Map(model, func(idx int, value declarations.Group) string {
		return value.Name
	})
}
