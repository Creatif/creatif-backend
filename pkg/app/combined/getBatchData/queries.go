package getBatchData

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type Variable struct {
	ID string `gorm:"primarykey"`

	Name      string `gorm:"index;uniqueIndex:unique_variable"`
	Behaviour string
	Groups    pq.StringArray `gorm:"type:text[]"`
	Metadata  datatypes.JSON
	Value     datatypes.JSON

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type QueriesMapVariable struct {
	ID string `gorm:"primarykey"`

	MapName   string
	Name      string `gorm:"index;uniqueIndex:unique_variable"`
	Behaviour string
	Groups    pq.StringArray `gorm:"type:text[]"`
	Metadata  datatypes.JSON
	Value     datatypes.JSON

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func queryVariableValue(variableIds []string) ([]Variable, error) {
	var variables []Variable
	if res := storage.Gorm().Raw(`SELECT n.id, n.name, n.behaviour, n.metadata, n.groups, n.created_at, n.updated_at, n.value FROM declarations.variables AS n WHERE n.id IN (?)
`, variableIds).Scan(&variables); res.Error != nil {
		return nil, res.Error
	}

	return variables, nil
}

func queryMapVariables(mapIds []string, model interface{}) error {
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
		(declarations.MapVariable{}).TableName(),
	)

	if res := storage.Gorm().Raw(sql, mapIds).Scan(model); res.Error != nil {
		return res.Error
	}

	return nil
}
