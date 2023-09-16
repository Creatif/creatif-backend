package getMap

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"strings"
)

func queryVariables(mapId string, fields []string, model interface{}) error {
	resolvedFields := strings.Join(sdk.Map(fields, func(idx int, value string) string {
		return fmt.Sprintf("n.%s", value)
	}), ",")

	delimiter := ""
	if len(resolvedFields) > 0 {
		delimiter = ","
	}

	sql := fmt.Sprintf(`
SELECT 
    n.id,
    n.name%s %s
		FROM %s AS n
		WHERE n.map_id = ?
`,
		delimiter,
		resolvedFields,
		(declarations.MapVariable{}).TableName(),
	)

	if res := storage.Gorm().Raw(sql, mapId).Scan(model); res.Error != nil {
		return res.Error
	}

	return nil
}

func queryMap(mapName string) (declarations.Map, error) {
	var m declarations.Map
	if err := storage.GetBy((declarations.Map{}).TableName(), "name", mapName, &m); err != nil {
		return declarations.Map{}, err
	}

	return m, nil
}
