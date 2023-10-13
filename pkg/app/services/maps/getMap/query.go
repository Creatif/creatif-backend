package getMap

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"strings"
)

func queryVariables(mapId, localeID string, fields []string, model interface{}) error {
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
		WHERE n.map_id = ? AND locale_id = ?
`,
		delimiter,
		resolvedFields,
		(declarations.MapVariable{}).TableName(),
	)

	if res := storage.Gorm().Raw(sql, mapId, localeID).Scan(model); res.Error != nil {
		return res.Error
	}

	return nil
}

func queryMap(projectId, mapName, localeID string) (declarations.Map, error) {
	var m declarations.Map
	if res := storage.Gorm().Where("name = ? AND project_id = ? AND locale_id = ?", mapName, projectId, localeID).First(&m); res.Error != nil {
		return declarations.Map{}, res.Error
	}

	return m, nil
}
