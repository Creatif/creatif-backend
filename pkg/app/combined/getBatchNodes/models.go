package getBatchNodes

import (
	"github.com/lib/pq"
	"time"
)

type node struct {
	Name string
	Type string
}

type GetBatchedNodesModel struct {
	Nodes []node
}

func NewGetBatchedNodesModel(nodes map[string]string) GetBatchedNodesModel {
	models := make([]node, len(nodes))
	for name, t := range nodes {
		models = append(models, node{
			Name: name,
			Type: t,
		})
	}

	return GetBatchedNodesModel{
		Nodes: models,
	}
}

type View struct {
	ID string `json:"id"`

	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Behaviour string         `json:"behaviour"`
	Groups    pq.StringArray `json:"groups"`
	Metadata  interface{}    `json:"metadata"`
	Value     interface{}    `json:"value"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model []NodeWithValueQuery) map[string]View {
	views := make(map[string]View)
	for _, value := range model {
		views[value.Name] = View{
			ID:        value.ID,
			Name:      value.Name,
			Type:      value.Type,
			Behaviour: value.Behaviour,
			Groups:    value.Groups,
			Metadata:  []byte(value.Metadata),
			Value:     []byte(value.Value),
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}
	}

	return views
}
