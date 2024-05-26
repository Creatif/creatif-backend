package connections

import (
	"creatif/pkg/app/domain/published"
	"fmt"
)

func getConnectionsSql() string {
	return fmt.Sprintf(`
SELECT 
    ref.child_id AS child,
    ref.parent_id AS parent
FROM %s AS ref
WHERE version_id = ? AND project_id = ? AND (ref.child_id IN(?) OR ref.parent_id IN(?))
`,
		(published.PublishedReference{}).TableName(),
	)
}
