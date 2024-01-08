package paginateLists

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/variables/paginateVariables/pagination"
	"time"
)

type LogicModel struct {
	variables      []Variable
	paginationInfo pagination.PaginationInfo
}

type View struct {
	ID      string `json:"id"`
	ShortID string `json:"shortId"`
	Name    string `json:"name"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(models []declarations.List) ([]View, error) {
	views := make([]View, len(models))
	for i, value := range models {
		views[i] = View{
			ID:        value.ID,
			Name:      value.Name,
			ShortID:   value.ShortID,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}
	}

	return views, nil
}
