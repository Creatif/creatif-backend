package paginateMaps

import (
	"creatif/pkg/app/domain/declarations"
	"time"
)

type View struct {
	ID      string `json:"id"`
	ShortID string `json:"shortId"`
	Name    string `json:"name"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(models []declarations.Map) ([]View, error) {
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
