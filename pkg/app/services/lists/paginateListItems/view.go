package paginateListItems

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/variables/paginateVariables/pagination"
	"creatif/pkg/lib/sdk"
	"time"
)

type LogicModel struct {
	variables      []Variable
	paginationInfo pagination.PaginationInfo
}

type View struct {
	ID        string      `json:"id"`
	Index     string      `json:"index"`
	ShortID   string      `json:"shortId"`
	Name      string      `json:"name"`
	Groups    []string    `json:"groups"`
	Behaviour string      `json:"behaviour"`
	Metadata  interface{} `json:"metadata"`
	Value     interface{} `json:"value"`

	CreatedAt time.Time `gorm:"<-:createProject" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(models []declarations.ListVariable) []View {
	return sdk.Map(models, func(idx int, value declarations.ListVariable) View {
		return View{
			ID:        value.ID,
			Name:      value.Name,
			Index:     value.Index,
			ShortID:   value.ShortID,
			Groups:    value.Groups,
			Value:     value.Value,
			Behaviour: value.Behaviour,
			Metadata:  value.Metadata,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}
	})
}
