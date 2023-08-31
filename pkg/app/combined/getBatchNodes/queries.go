package getBatchNodes

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type Node struct {
	ID uuid.UUID `gorm:"primarykey"`

	Name      string `gorm:"index;uniqueIndex:unique_node"`
	Type      string
	Behaviour string
	Groups    pq.StringArray `gorm:"type:text[]"`
	Metadata  datatypes.JSON
	Value     datatypes.JSON

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func queryValue(nodeIds []uuid.UUID) ([]Node, error) {
	var nodes []Node
	if res := storage.Gorm().Raw(fmt.Sprintf(`SELECT n.id, n.name, n.type, n.behaviour, n.metadata, n.groups, n.created_at, n.updated_at, vn.value FROM declarations.nodes AS n
INNER JOIN %s AS an ON n.id = an.declaration_node_id
INNER JOIN %s AS vn ON an.id = vn.assignment_node_id
WHERE n.id IN (?)
`, (assignments.Node{}).TableName(), (assignments.ValueNode{}).TableName()), nodeIds).Scan(&nodes); res.Error != nil {
		return nil, res.Error
	}

	return nodes, nil
}
