package addGroups

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/lib/sdk"
)

func newView(model []app.Group) []string {
	return sdk.Map(model, func(idx int, value app.Group) string {
		return value.Name
	})
}
