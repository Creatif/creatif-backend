package paginateProjects

import (
	"creatif/pkg/lib/sdk"
	"time"
)

type View struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`

	MapsNum  int `json:"mapsNum"`
	ListsNum int `json:"listsNum"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(models []QueryModel) []View {
	return sdk.Map(models, func(idx int, value QueryModel) View {
		return View{
			ID:        value.ID,
			Name:      value.Name,
			State:     value.State,
			MapsNum:   value.MapsNum,
			ListsNum:  value.ListsNum,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}
	})
}
