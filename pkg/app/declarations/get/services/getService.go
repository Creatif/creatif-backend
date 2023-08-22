package services

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type Node struct {
	ID string `gorm:"primarykey"`

	Name      string         `gorm:"index;uniqueIndex:unique_node"`
	Type      string         // text,image,file,boolean
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
	if res := storage.Gorm().Raw(`
SELECT n.id, n.name, n.type, n.behaviour, n.metadata, n.groups, n.created_at, n.updated_at, vn.value FROM declarations.nodes AS n
	INNER JOIN assignments.nodes AS an ON n.id = an.declaration_node_id
	INNER JOIN assignments.value_node AS vn ON an.id = vn.assignment_node_id
	WHERE n.id = ?
`, nodeId).Scan(&node); res.Error != nil {
		return Node{}, res.Error
	}

	return node, nil
}
