package pagination

type Cursor struct {
	next      *string
	prev      *string
	direction string // forward or backwards
}

type Rule struct {
	field   string
	orderBy string
}

type Pagination struct {
	sql    string
	rules  []Rule
	cursor Cursor
}

func NewCursor(next string, prev string) Cursor {
	return Cursor{
		next: &next,
		prev: &prev,
	}
}

func NewPagination(sql string, rules []Rule, cursor Cursor) Pagination {
	return Pagination{
		sql:    sql,
		rules:  rules,
		cursor: cursor,
	}
}
