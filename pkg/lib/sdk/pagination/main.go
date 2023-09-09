package pagination

import (
	"creatif/pkg/lib/storage"
	"fmt"
	"strings"
)

const ASC = "ASC"
const DESC = "DESC"

const DIRECTION_FORWARD = "forward"
const DIRECTION_BACKWARDS = "backwards"

type SortableIDModel struct {
	ID string `gorm:"primarykey;type:text CHECK(length(id)=27)"`
}

type orderByRule struct {
	field   string
	orderBy string
}

func NewOrderByRule(field, orderBy string) orderByRule {
	return orderByRule{
		field:   field,
		orderBy: strings.ToUpper(orderBy),
	}
}

type Pagination struct {
	table     string
	sql       string
	limit     int
	cursor    string
	rule      orderByRule
	direction string
}

func NewPagination(table, sql string, rules orderByRule, cursor string, limit int, direction string) *Pagination {
	return &Pagination{
		table:     table,
		sql:       sql,
		rule:      rules,
		limit:     limit,
		direction: direction,
		cursor:    cursor,
	}
}

func (p Pagination) Paginate(model interface{}) error {
	isFirstPage := p.cursor == ""
	if isFirstPage {
		id, err := getInitialID(p.table, p.rule.orderBy)
		if err != nil {
			return err
		}

		cur, err := NewCursor(p.cursor, id, p.rule.field, p.rule.orderBy)
		if err != nil {
			return err
		}

		operator := getInitialOperator(DIRECTION_FORWARD, p.rule.orderBy)
		if res := storage.Gorm().Raw(fmt.Sprintf("%s WHERE id %s '%s' ORDER BY %s %s LIMIT %d", p.sql, operator, cur.ID, cur.Field, cur.OrderBy, p.limit)).Scan(model); res.Error != nil {
			return res.Error
		}

		return nil
	}

	var cur *cursor
	cur, err := NewCursor(p.cursor, "", p.rule.field, p.rule.orderBy)
	if err != nil {
		return err
	}

	operator := getOperator(p.direction, cur.OrderBy)
	if res := storage.Gorm().Raw(fmt.Sprintf("%s WHERE id %s '%s' ORDER BY %s %s LIMIT %d", p.sql, operator, cur.ID, cur.Field, cur.OrderBy, p.limit)).Scan(model); res.Error != nil {
		return cursor{}, res.Error
	}

	return cur, nil
}

func (p Pagination) PaginationInfo(prevCur cursor, nextId string) (PaginationInfo, error) {
	nextCursor, err := resolveCursor(nextId, prevCur.Field, prevCur.OrderBy)
	if err != nil {
		return PaginationInfo{}, err
	}

	prev, err := encodeCursor(prevCur)
	if err != nil {
		return PaginationInfo{}, err
	}

	return PaginationInfo{
		Next:    nextCursor,
		Prev:    prev,
		NextURL: "",
		PrevURL: "",
	}, nil
}
