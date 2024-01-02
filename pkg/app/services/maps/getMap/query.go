package getMap

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"strings"
)

func queryVariables(mapId string, fields, groups []string, model interface{}) error {
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
		WHERE n.map_id = ?
%s
`,
		delimiter,
		resolvedFields,
		(declarations.MapVariable{}).TableName(),
		groupsWhereClause,
	)

	if res := storage.Gorm().Raw(sql, mapId).Scan(model); res.Error != nil {
		return res.Error
	}

	return nil
}

func queryMap(projectId, name string) (declarations.Map, error) {
	var m declarations.Map
	if res := storage.Gorm().Where(fmt.Sprintf("project_id = ? AND (name = ? OR id = ? OR short_id = ?)"), projectId, name, name, name).First(&m); res.Error != nil {
		return declarations.Map{}, res.Error
	}

	return m, nil
}
