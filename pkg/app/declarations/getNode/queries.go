package getNode

import (
	"creatif/pkg/app/domain/assignments"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"strings"
	"time"
)

type NodeWithValueQuery struct {
	ID string

	Name      string
	Behaviour string // readonly,modifiable
	Groups    pq.StringArray
	Metadata  datatypes.JSON
	Value     interface{}

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func queryValue(name string, fields []string) (Node, error) {
	resolvedFields := strings.Join(sdk.Map(fields, func(idx int, value string) string {
		if value != "value" {
			return fmt.Sprintf("n.%s", value)
		}

		return fmt.Sprintf("vn.value")
	}), ",")

	var node Node
	if res := storage.Gorm().Raw(fmt.Sprintf(`SELECT n.id, n.name, %s FROM declarations.nodes AS n
INNER JOIN %s AS an ON n.id = an.declaration_node_id
INNER JOIN %s AS vn ON an.id = vn.assignment_node_id
WHERE n.name = ?
`, resolvedFields, (assignments.Node{}).TableName(), (assignments.ValueNode{}).TableName()), name).Scan(&node); res.Error != nil {
		return Node{}, res.Error
	}

	return node, nil
}
