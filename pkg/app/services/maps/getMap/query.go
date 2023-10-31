package getMap

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"strings"
)

func queryVariables(mapId, localeID string, fields, groups []string, model interface{}) error {
	resolvedFields := strings.Join(sdk.Map(fields, func(idx int, value string) string {
		return fmt.Sprintf("n.%s", value)
	}), ",")

	delimiter := ""
	if len(resolvedFields) > 0 {
		delimiter = ","
	}

	var groupsWhereClause string
	if len(groups) != 0 {
		groupsWhereClause = fmt.Sprintf("AND '{%s}'::text[] && %s", strings.Join(groups, ","), "n.groups")
	}

	sql := fmt.Sprintf(`
SELECT 
    n.id,
    n.name%s %s
		FROM %s AS n
		WHERE n.map_id = ? AND n.locale_id = ?
%s
`,
		delimiter,
		resolvedFields,
		(declarations.MapVariable{}).TableName(),
		groupsWhereClause,
	)

	if res := storage.Gorm().Raw(sql, mapId, localeID).Scan(model); res.Error != nil {
		return res.Error
	}

	return nil
}

func queryMap(projectId, mapId, val, localeID string) (declarations.Map, error) {
	var m declarations.Map
	if res := storage.Gorm().Where(fmt.Sprintf("%s AND project_id = ? AND locale_id = ?", mapId), val, projectId, localeID).First(&m); res.Error != nil {
		return declarations.Map{}, res.Error
	}

	return m, nil
}
