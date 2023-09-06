package pagination

const ASC = "ASC"
const DESC = "DESC"

const DIRECTION_FORWARD = "forward"
const DIRECTION_BACKWARDS = "backwards"

type orderByRule struct {
	field   string
	orderBy string
}

type Pagination struct {
	table  string
	sql    string
	cursor string
	rule   orderByRule
}

func NewPagination(table, sql string, rules orderByRule, cursor string) Pagination {
	return Pagination{
		table:  table,
		sql:    sql,
		rule:   rules,
		cursor: cursor,
	}
}

func (p Pagination) Create() (string, error) {

	return "", nil
}
