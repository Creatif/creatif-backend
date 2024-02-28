package publish

import "fmt"

func buildInsertSql(table string, params []string) string {
	return fmt.Sprintf(`
	INSERT INTO %s (%s) VALUES (%s)
`)
}
