package getMap

import (
	"creatif/pkg/app/domain/declarations"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

var validFields = []string{
	"behaviour",
	"metadata",
	"groups",
	"created_at",
	"updated_at",
}

type Model struct {
	Name   string
	Fields []string

	validFields []string
	// TODO: Add project ID prop here
}

func NewModel(name string, fields []string) Model {
	return Model{
		Name:        name,
		Fields:      fields,
		validFields: validFields,
	}
}

type LogicModel struct {
	nodeMap declarations.Map
	nodes   []Node
}

type NamesOnlyView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Node struct {
	ID string `json:"id" gorm:"primarykey"`

	Name      string         `json:"name" gorm:"index;uniqueIndex:unique_node"`
	Value     datatypes.JSON `json:"value"`
	Behaviour string         `json:"behaviour"`
	Groups    pq.StringArray `json:"groups" gorm:"type:text[]"`
	Metadata  datatypes.JSON `json:"metadata"`

	CreatedAt time.Time `json:"createdAt" gorm:"<-:create"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type View struct {
	ID    string                   `json:"id"`
	Name  string                   `json:"name"`
	Nodes []map[string]interface{} `json:"nodes"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model LogicModel, returnFields []string) View {
	m := make([]map[string]interface{}, 0)

	for _, n := range model.nodes {
		o := make(map[string]interface{})

		o["id"] = n.ID
		o["name"] = n.Name

		for _, f := range returnFields {
			if f == "groups" {
				o["groups"] = n.Groups
			}

			if f == "value" {
				o["value"] = n.Value
			}

			if f == "behaviour" {
				o["behaviour"] = n.Behaviour
			}

			if f == "metadata" {
				o["metadata"] = n.Metadata
			}

			if f == "value" {
				o["value"] = n.Value
			}

			if f == "created_at" {
				o["createdAt"] = n.CreatedAt
			}

			if f == "updated_at" {
				o["updatedAt"] = n.UpdatedAt
			}
		}

		m = append(m, o)
	}

	return View{
		ID:        model.nodeMap.ID,
		Name:      model.nodeMap.Name,
		Nodes:     m,
		CreatedAt: model.nodeMap.CreatedAt,
		UpdatedAt: model.nodeMap.UpdatedAt,
	}
}
