package paginateProjects

import (
	"creatif/pkg/app/services/variables/paginateVariables/pagination"
	"creatif/pkg/lib/sdk"
	"time"
)

type LogicModel struct {
	variables      []Variable
	paginationInfo pagination.PaginationInfo
}

type View struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`

	VariableNum int `json:"variableNum"`
	MapsNum     int `json:"mapsNum"`
	ListsNum    int `json:"listsNum"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(models []QueryModel) []View {
	return sdk.Map(models, func(idx int, value QueryModel) View {
		return View{
			ID:          value.ID,
			Name:        value.Name,
			State:       value.State,
			VariableNum: value.VariableNum,
			MapsNum:     value.MapsNum,
			ListsNum:    value.ListsNum,
			CreatedAt:   value.CreatedAt,
			UpdatedAt:   value.UpdatedAt,
		}
	})
}
