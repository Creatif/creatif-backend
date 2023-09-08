package getBatchNodes

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type Node struct {
	ID string `gorm:"primarykey"`

	Name      string `gorm:"index;uniqueIndex:unique_node"`
	Behaviour string
	Groups    pq.StringArray `gorm:"type:text[]"`
	Metadata  datatypes.JSON
	Value     datatypes.JSON

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MapNode struct {
	ID string `gorm:"primarykey"`

	MapName   string
	Name      string `gorm:"index;uniqueIndex:unique_node"`
	Behaviour string
	Groups    pq.StringArray `gorm:"type:text[]"`
	Metadata  datatypes.JSON
	Value     datatypes.JSON

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func queryNodesValue(nodeIds []string) ([]Node, error) {
	var nodes []Node
	if res := storage.Gorm().Raw(fmt.Sprintf(`SELECT n.id, n.name, n.behaviour, n.metadata, n.groups, n.created_at, n.updated_at, vn.value FROM declarations.nodes AS n
INNER JOIN %s AS an ON n.id = an.declaration_node_id
INNER JOIN %s AS vn ON an.id = vn.assignment_node_id
WHERE n.id IN (?)
`, (assignments.Node{}).TableName(), (assignments.ValueNode{}).TableName()), nodeIds).Scan(&nodes); res.Error != nil {
		return nil, res.Error
	}

	return nodes, nil
}

func queryMapValues(mapIds []string) ([]MapNode, error) {
	sql := fmt.Sprintf(`SELECT 
    n.id, 
    m.name AS mapName,
    n.name, 
    n.groups, 
    n.behaviour, 
    vn.value,
    n.metadata, 
    n.created_at, 
    n.updated_at
		FROM %s AS mn
		INNER JOIN %s AS n ON n.id = mn.node_id
		INNER JOIN %s AS m ON m.id = mn.map_id
		INNER JOIN %s AS an ON an.declaration_node_id = n.id
		INNER JOIN %s AS vn ON vn.assignment_node_id = an.id
		WHERE m.id IN (?)
`,
		(declarations.MapNode{}).TableName(),
		(declarations.Node{}).TableName(),
		(declarations.Map{}).TableName(),
		(assignments.Node{}).TableName(),
		(assignments.ValueNode{}).TableName(),
	)

	var nodes []MapNode
	if res := storage.Gorm().Raw(sql, mapIds).Scan(&nodes); res.Error != nil {
		return nil, res.Error
	}

	return nodes, nil
}
