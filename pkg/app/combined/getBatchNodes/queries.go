package getBatchNodes

import (
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

type QueriesMapNode struct {
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
	if res := storage.Gorm().Raw(`SELECT n.id, n.name, n.behaviour, n.metadata, n.groups, n.created_at, n.updated_at, n.value FROM declarations.nodes AS n WHERE n.id IN (?)
`, nodeIds).Scan(&nodes); res.Error != nil {
		return nil, res.Error
	}

	return nodes, nil
}

func queryMapNodes(mapIds []string, model interface{}) error {
	sql := fmt.Sprintf(`
SELECT 
    n.id,
    n.name,
    m.name AS mapName,
    n.groups,
    n.metadata,
    n.value,
    n.created_at,
    n.updated_at
		FROM %s AS m
		INNER JOIN %s AS n ON m.id = n.map_id
		WHERE m.id IN(?)
`,
		(declarations.Map{}).TableName(),
		(declarations.MapNode{}).TableName(),
	)

	if res := storage.Gorm().Raw(sql, mapIds).Scan(model); res.Error != nil {
		return res.Error
	}

	return nil
}
