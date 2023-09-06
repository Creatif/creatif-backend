package pagination

import "fmt"

func createInitialQuery(table, direction string) string {
	return fmt.Sprintf(`SELECT id, created_at FROM %s ORDER BY created_at %s LIMIT 1`, table, direction)
}

func createPaginationQuery(sql, direction, id, createdAt string, rule orderByRule) string {
	querySegment := "<="
	if direction == ASC {
		querySegment = ">="
	}

	return fmt.Sprintf(`%s WHERE id %s %s AND created_at %s %s ORDER BY %s %s`, sql, querySegment, id, querySegment, createdAt, rule.field, rule.orderBy)
}
