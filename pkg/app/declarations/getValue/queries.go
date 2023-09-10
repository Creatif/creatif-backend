package getValue

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/lib/storage"
	"fmt"
)

func queryValue(name string) (Node, error) {
	var node Node
	if res := storage.Gorm().Raw(fmt.Sprintf(`SELECT vn.value FROM declarations.nodes AS n
INNER JOIN %s AS an ON n.id = an.declaration_node_id
INNER JOIN %s AS vn ON an.id = vn.assignment_node_id
WHERE n.name = ?
`, (assignments.Node{}).TableName(), (assignments.ValueNode{}).TableName()), name).Scan(&node); res.Error != nil {
		return Node{}, res.Error
	}

	return node, nil
}
