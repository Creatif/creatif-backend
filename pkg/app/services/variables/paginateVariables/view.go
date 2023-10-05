package paginateVariables

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
	ID        string      `json:"id"`
	ProjectID string      `json:"projectID"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Groups    []string    `json:"groups"`
	Behaviour string      `json:"behaviour"`
	Metadata  interface{} `json:"metadata"`
	Value     interface{} `json:"value"`

	CreatedAt time.Time `gorm:"<-:createProject" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ViewParameters struct {
	Field        string   `json:"field"`
	OrderBy      string   `json:"orderBy"`
	Groups       []string `json:"groups"`
	Limit        int      `json:"limit"`
	PaginationID string   `json:"paginationId"`
}

type ViewPaginationInfo struct {
	Next       string         `json:"next"`
	Prev       string         `json:"prev"`
	Parameters ViewParameters `json:"parameters"`
}

type PaginatedView struct {
	Items          []View             `json:"items"`
	PaginationInfo ViewPaginationInfo `json:"paginationInfo"`
}

func newView(model LogicModel) PaginatedView {
	views := sdk.Map(model.variables, func(idx int, value Variable) View {
		return View{
			ID:        value.ID,
			ProjectID: value.ProjectID,
			Name:      value.Name,
			Groups:    value.Groups,
			Value:     value.Value,
			Behaviour: value.Behaviour,
			Metadata:  value.Metadata,
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}
	})

	p := PaginatedView{
		Items: views,
		PaginationInfo: ViewPaginationInfo{
			Next: model.paginationInfo.Next,
			Prev: model.paginationInfo.Prev,
			Parameters: ViewParameters{
				PaginationID: model.paginationInfo.Parameters.PaginationID,
				Field:        model.paginationInfo.Parameters.Field,
				OrderBy:      model.paginationInfo.Parameters.OrderBy,
				Groups:       model.paginationInfo.Parameters.Groups,
				Limit:        model.paginationInfo.Parameters.Limit,
			},
		},
	}

	return p
}
