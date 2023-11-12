package paginateProjects

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/services/variables/paginateVariables/pagination"
	"creatif/pkg/lib/sdk"
	"time"
)

type LogicModel struct {
	variables      []Variable
	paginationInfo pagination.PaginationInfo
}

type View struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	APIKey string `json:"apiKey"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(models []app.Project) []View {
	return sdk.Map(models, func(idx int, value app.Project) View {
		return View{
			ID:        value.ID,
			APIKey:    value.APIKey,
			Name:      value.Name,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}
	})
}
