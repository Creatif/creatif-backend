package deleteList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"fmt"
)

func removeConnections(projectId, structureId string) error {
	sql := fmt.Sprintf(`
DELETE FROM %s AS c
USING declarations.lists AS l, declarations.list_variables AS lv
WHERE l.project_id = ? AND l.id = ? AND c.parent_variable_id = lv.id AND c.child_variable_id = lv.id
`, (declarations.Connection{}).TableName())

	if res := storage.Gorm().Exec(sql, projectId, structureId); res.Error != nil {
		return res.Error
	}

	return nil
}
