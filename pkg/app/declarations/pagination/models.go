package pagination

import (
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/sdk/pagination"
	"time"
)

type PaginationModel struct {
	// this can be project name
	WithValue bool
	SortField string
	Limit     int
	SortOrder string
	Cursor    string
	// TODO: Add project ID prop here
}

func NewModel(withValue bool, sortField, sortOrder string, limit int) PaginationModel {
	return PaginationModel{
		WithValue: withValue,
		SortField: sortField,
		Limit:     limit,
		SortOrder: sortOrder,
	}
}

type LogicModelWithoutValue struct {
	nodes          []NodeWithoutValue
	paginationInfo pagination.PaginationInfo
}

type LogicModelWithValue struct {
	nodes          []NodeWithValue
	paginationInfo pagination.PaginationInfo
}

type View struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Groups    []string               `json:"groups"`
	Behaviour string                 `json:"behaviour"`
	Metadata  map[string]interface{} `json:"metadata"`
	Value     interface{}            `json:"value,omitempty"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ViewPaginationInfo struct {
	Next    string `json:"next"`
	Prev    string `json:"prev"`
	NextURL string `json:"nextURL"`
	PrevURL string `json:"prevURL"`
}

type PaginatedView struct {
	Items          []View             `json:"items"`
	PaginationInfo ViewPaginationInfo `json:"paginationInfo"`
}

func newView(model interface{}) PaginatedView {
	if m, ok := model.(LogicModelWithoutValue); ok {
		views := sdk.Map(m.nodes, func(idx int, value NodeWithoutValue) View {
			return View{
				ID:        value.ID,
				Name:      value.Name,
				Groups:    value.Groups,
				Behaviour: value.Behaviour,
				Metadata:  sdk.UnmarshalToMap([]byte(value.Metadata)),
				CreatedAt: value.CreatedAt,
				UpdatedAt: value.UpdatedAt,
			}
		})

		return PaginatedView{
			Items: views,
			PaginationInfo: ViewPaginationInfo{
				Next:    m.paginationInfo.Next,
				Prev:    m.paginationInfo.Prev,
				NextURL: m.paginationInfo.NextURL,
				PrevURL: m.paginationInfo.PrevURL,
			},
		}
	}

	m := model.(LogicModelWithValue)
	views := sdk.Map(m.nodes, func(idx int, value NodeWithValue) View {
		return View{
			ID:        value.ID,
			Name:      value.Name,
			Groups:    value.Groups,
			Value:     value.Value,
			Behaviour: value.Behaviour,
			Metadata:  sdk.UnmarshalToMap([]byte(value.Metadata)),
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}
	})

	return PaginatedView{
		Items: views,
		PaginationInfo: ViewPaginationInfo{
			Next:    m.paginationInfo.Next,
			Prev:    m.paginationInfo.Prev,
			NextURL: m.paginationInfo.NextURL,
			PrevURL: m.paginationInfo.PrevURL,
		},
	}
}
