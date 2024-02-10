package addGroups

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
)

type View struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func newView(model []declarations.Group) []View {
	return sdk.Map(model, func(idx int, value declarations.Group) View {
		return View{
			ID:   value.ID,
			Name: value.Name,
		}
	})
}
