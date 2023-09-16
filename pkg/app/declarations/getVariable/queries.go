package getVariable

import (
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"strings"
	"time"
)

type VariableWithValueQuery struct {
	ID string

	Name      string
	Behaviour string // readonly,modifiable
	Groups    pq.StringArray
	Metadata  datatypes.JSON
	Value     interface{}

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func queryValue(name string, fields []string) (Variable, error) {
	resolvedFields := strings.Join(sdk.Map(fields, func(idx int, value string) string {
		return fmt.Sprintf("n.%s", value)
	}), ",")

	var variable Variable
	if res := storage.Gorm().Raw(fmt.Sprintf(`SELECT n.id, n.name, %s FROM declarations.variables AS n
WHERE n.name = ?
`, resolvedFields), name).Scan(&variable); res.Error != nil {
		return Variable{}, res.Error
	}

	return variable, nil
}
