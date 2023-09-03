package get

import (
	"creatif/pkg/lib/sdk"
	"github.com/google/uuid"
	"time"
)

type PaginationModel struct {
	// this can be project name
	WithValue bool
	SortField string
	Limit     int
	SortOrder string
	// TODO: Add project ID prop here
}

func NewPaginationModel(withValue bool, sortField, sortOrder string, limit int) PaginationModel {
	return PaginationModel{
		WithValue: withValue,
		SortField: sortField,
		Limit:     limit,
		SortOrder: sortOrder,
	}
}

type ViewWithoutValue struct {
	ID        uuid.UUID              `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Groups    []string               `json:"groups"`
	Behaviour string                 `json:"behaviour"`
	Metadata  map[string]interface{} `json:"metadata"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type View struct {
	ID        uuid.UUID              `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Groups    []string               `json:"groups"`
	Behaviour string                 `json:"behaviour"`
	Metadata  map[string]interface{} `json:"metadata"`
	Value     interface{}            `json:"value"`

	CreatedAt time.Time `gorm:"<-:create" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model NodeWithValueQuery) View {
	return View{
		ID:        model.ID,
		Name:      model.Name,
		Value:     model.Value,
		Groups:    model.Groups,
		Behaviour: model.Behaviour,
		Metadata:  sdk.UnmarshalToMap([]byte(model.Metadata)),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
