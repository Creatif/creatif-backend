package getMap

import (
	"creatif/pkg/app/domain/declarations"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

var validFields = []string{
	"id",
	"name",
	"type",
	"behaviour",
	"metadata",
	"groups",
	"created_at",
	"updated_at",
}

type GetMapModel struct {
	// this can be map name or an id of the map
	ID string
	// this can be, 'full' | names
	Return string
	// this can be individual fields of the node to return, reduces returned data
	// if the user needs only metadata, only metadata will be returned
	// name is always returned
	Fields []string

	validFields []string
	// TODO: Add project ID prop here
}

func NewGetMapModel(id string, ret string, fields []string) GetMapModel {
	return GetMapModel{
		ID:          id,
		Return:      ret,
		Fields:      fields,
		validFields: validFields,
	}
}

type NamesOnlyView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FullNode struct {
	ID string `json:"id" gorm:"primarykey"`

	Name      string         `json:"name" gorm:"index;uniqueIndex:unique_node"`
	Type      string         `json:"type"`
	Value     interface{}    `json:"value"`
	Behaviour string         `json:"behaviour"`
	Groups    pq.StringArray `json:"groups" gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata  datatypes.JSON `json:"metadata"`

	CreatedAt time.Time `json:"createdAt" gorm:"<-:create"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CustomNode struct {
	ID string `json:"id"`

	Name      string         `json:"name"`
	Type      string         `json:"type,omitempty"`
	Behaviour string         `json:"behaviour,omitempty"`
	Value     datatypes.JSON `json:"value"`
	Groups    pq.StringArray `json:"groups,omitempty" gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata  datatypes.JSON `json:"metadata,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"<-:create"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type View struct {
	ID    uuid.UUID   `json:"id"`
	Name  string      `json:"name"`
	Nodes interface{} `json:"nodes"`
}

func newView(model LogicModel) View {
	view := View{
		ID:   model.nodeMap.ID,
		Name: model.nodeMap.Name,
	}

	if model.strategy == "namesOnlyStrategy" {
		views := make([]NamesOnlyView, 0)
		for _, n := range model.nodes {
			views = append(views, NamesOnlyView{
				ID:   n.ID,
				Name: n.Name,
			})
		}

		view.Nodes = views
		return view
	}

	if model.strategy == "customFieldsStrategy" {
		views := make([]CustomNode, 0)
		for _, n := range model.nodes {
			views = append(views, CustomNode{
				ID:        n.ID,
				Name:      n.Name,
				Type:      n.Type,
				Behaviour: n.Behaviour,
				Groups:    n.Groups,
				Metadata:  n.Metadata,
				CreatedAt: n.CreatedAt,
				UpdatedAt: n.UpdatedAt,
			})
		}

		view.Nodes = views
		return view
	}

	view.Nodes = model.nodes
	return view
}

type LogicModel struct {
	nodeMap  declarations.Map
	nodes    []FullNode
	strategy string
}
