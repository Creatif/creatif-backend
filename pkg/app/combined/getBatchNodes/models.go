package getBatchNodes

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type node struct {
	Name string
	Type string
}

type GetBatchedNodesModel struct {
	Nodes []node

	mapIds  []uuid.UUID
	nodeIds []uuid.UUID
}

func NewGetBatchedNodesModel(nodes map[string]string) *GetBatchedNodesModel {
	models := make([]node, len(nodes))
	count := 0
	for name, t := range nodes {
		models[count] = node{
			Name: name,
			Type: t,
		}
		count++
	}

	return &GetBatchedNodesModel{
		Nodes: models,
	}
}

type View struct {
	ID uuid.UUID `json:"id"`

	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Behaviour string         `json:"behaviour"`
	Groups    pq.StringArray `json:"groups"`
	Metadata  interface{}    `json:"metadata"`
	Value     interface{}    `json:"value"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func newView(model []Node) map[string]View {
	nodes := make(map[string]View)
	for _, m := range model {
		nodes[m.Name] = View{
			ID:        m.ID,
			Name:      m.Name,
			Type:      m.Type,
			Behaviour: m.Behaviour,
			Groups:    m.Groups,
			Metadata:  m.Metadata,
			Value:     m.Value,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}

	return nodes
}

func (a *GetBatchedNodesModel) Validate() map[string]string {
	v := map[string]interface{}{
		"validNodes": a.Nodes,
		"validNames": a.Nodes,
	}

	if err := validation.Validate(v,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("validNodes", validation.By(func(value interface{}) error {
				nodes := value.([]node)

				for _, t := range nodes {
					if t.Name == "" {
						return errors.New("Node name cannot be empty")
					}
				}

				for _, t := range nodes {
					if t.Type != "node" && t.Type != "map" {
						return errors.New(fmt.Sprintf("Invalid type in node with name '%s'. Valid type are 'map' and 'node'", t.Name))
					}
				}

				return nil
			})),
			validation.Key("validNames", validation.By(func(value interface{}) error {
				nodes := value.([]node)

				nodeNames := sdk.Filter(sdk.Map(nodes, func(idx int, value node) string {
					if value.Type == "node" {
						return value.Name
					}

					return ""
				}), func(idx int, value string) bool {
					return value != ""
				})

				mapNames := sdk.Filter(sdk.Map(nodes, func(idx int, value node) string {
					if value.Type == "map" {
						return value.Name
					}

					return ""
				}), func(idx int, value string) bool {
					return value != ""
				})

				var foundNodes []declarations.Node
				if res := storage.Gorm().Table((declarations.Node{}).TableName()).Select("ID").Where("name IN (?)", nodeNames).Find(&foundNodes); res.Error != nil {
					return errors.New("One of the node or map names given is invalid or does not exist.")
				}

				var maps []declarations.Map
				if res := storage.Gorm().Table((declarations.Map{}).TableName()).Select("ID").Where("name IN (?)", mapNames).Find(&maps); res.Error != nil {
					return errors.New("One of the node or map names given is invalid or does not exist.")
				}

				if (len(nodeNames) + len(mapNames)) != len(nodes) {
					return errors.New("One of the node or map names given is invalid or does not exist.")
				}

				a.nodeIds = sdk.Map(foundNodes, func(idx int, value declarations.Node) uuid.UUID {
					return value.ID
				})

				a.mapIds = sdk.Map(maps, func(idx int, value declarations.Map) uuid.UUID {
					return value.ID
				})

				return nil
			})),
		),
	); err != nil {
		var e map[string]string
		b, err := json.Marshal(err)
		if err != nil {
			return map[string]string{
				"unrecoverable": "An internal validation error occurred. This should not happen. Please, submit a bug.",
			}
		}

		if err := json.Unmarshal(b, &e); err != nil {
			return map[string]string{
				"unrecoverable": "An internal validation error occurred. This should not happen. Please, submit a bug.",
			}
		}

		return e
	}

	return nil
}
