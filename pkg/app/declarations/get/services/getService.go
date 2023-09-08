package services

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type Node struct {
	ID string `gorm:"primarykey"`

	Name      string         `gorm:"index;uniqueIndex:unique_node"`
	Behaviour string         // readonly,modifiable
	Groups    pq.StringArray `gorm:"type:text[]"` // if groups is set, group should be invalidated
	Metadata  datatypes.JSON
	Value     datatypes.JSON

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetService struct {
	id string
}

func NewGetService(id string) GetService {
	return GetService{id: id}
}

func (g GetService) GetNode(byId func(id string) (declarations.Node, error), byName func(name string) (declarations.Node, error)) (Node, error) {
	if sdk.IsValidUuid(g.id) {
		node, err := byId(g.id)
		if err != nil {
			return Node{}, err
		}

		return queryValue(node.ID)
	}

	node, err := byName(g.id)
	if err != nil {
		return Node{}, err
	}

	return queryValue(node.ID)
}

func queryValue(nodeId string) (Node, error) {
	var node Node
	if res := storage.Gorm().Raw(fmt.Sprintf(`SELECT n.id, n.name, n.behaviour, n.metadata, n.groups, n.created_at, n.updated_at, vn.value FROM declarations.nodes AS n
INNER JOIN %s AS an ON n.id = an.declaration_node_id
INNER JOIN %s AS vn ON an.id = vn.assignment_node_id
WHERE n.id = ?
`, (assignments.Node{}).TableName(), (assignments.ValueNode{}).TableName()), nodeId).Scan(&node); res.Error != nil {
		return Node{}, res.Error
	}

	return node, nil
}
