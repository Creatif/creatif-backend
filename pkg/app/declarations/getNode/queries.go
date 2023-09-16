package getNode

import (
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
		return fmt.Sprintf("n.%s", value)
	}), ",")

	var node Node
	if res := storage.Gorm().Raw(fmt.Sprintf(`SELECT n.id, n.name, %s FROM declarations.nodes AS n
WHERE n.name = ?
`, resolvedFields), name).Scan(&node); res.Error != nil {
		return Node{}, res.Error
	}

	return node, nil
}
