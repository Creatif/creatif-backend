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
	field      string
	orderBy    string
	groups     []string
	groupField string
}

func NewOrderByRule(field, orderBy, groupField string, groups []string) orderByRule {
	return orderByRule{
		field:      field,
		groupField: groupField,
		groups:     groups,
		orderBy:    strings.ToUpper(orderBy),
	}
}

type Pagination struct {
	table     string
	sql       string
	nextId    string
	prevId    string
	limit     int
	rule      orderByRule
	direction string
}

func NewPagination(table, sql string, rules orderByRule, nextId, prevId, direction string, limit int) *Pagination {
	return &Pagination{
		table:     table,
		sql:       sql,
		rule:      rules,
		nextId:    nextId,
		prevId:    prevId,
		direction: direction,
		limit:     limit,
	}
}

func (p Pagination) Paginate(model interface{}) error {
	isFirstPage := p.nextId == "" && p.prevId == ""
	if isFirstPage {
		id, err := getInitialID(p.table, p.rule.orderBy)
		if err != nil {
			return err
		}

		operator := getInitialOperator(DIRECTION_FORWARD, p.rule.orderBy)
		if len(p.rule.groups) > 0 {
			groups := strings.Join(p.rule.groups, ",")
			if res := storage.Gorm().Raw(fmt.Sprintf("%s WHERE id %s '%s' AND '{%s}'::text[] && %s ORDER BY %s %s LIMIT %d", p.sql, operator, id, groups, p.rule.groupField, p.rule.field, p.rule.orderBy, p.limit)).Scan(model); res.Error != nil {
				return res.Error
			}
		} else {
			if res := storage.Gorm().Raw(fmt.Sprintf("%s WHERE id %s '%s' ORDER BY %s %s LIMIT %d", p.sql, operator, id, p.rule.field, p.rule.orderBy, p.limit)).Scan(model); res.Error != nil {
				return res.Error
			}
		}

		return nil
	}

	operator := getOperator(p.direction, p.rule.orderBy)
	if res := storage.Gorm().Raw(fmt.Sprintf("%s WHERE id %s '%s' ORDER BY %s %s LIMIT %d", p.sql, operator, p.nextId, p.rule.field, p.rule.orderBy, p.limit)).Scan(model); res.Error != nil {
		return res.Error
	}

	return nil
}

func (p Pagination) PaginationInfo(nextId, prevId string) (PaginationInfo, error) {
	var nextCursor string
	if nextId != "" {
		n, err := encodeCursor(CursorFromData(nextId, prevId, p.rule.field, p.rule.orderBy, DIRECTION_FORWARD, p.limit))
		if err != nil {
			return PaginationInfo{}, err
		}

		nextCursor = n
	}

	var prevCursor string
	if prevId != "" {
		n, err := encodeCursor(CursorFromData(nextId, prevId, p.rule.field, p.rule.orderBy, DIRECTION_BACKWARDS, p.limit))
		if err != nil {
			return PaginationInfo{}, err
		}
		
		prevCursor = n
	}

	return PaginationInfo{
		Next:    nextCursor,
		Prev:    prevCursor,
		NextURL: "",
		PrevURL: "",
	}, nil
}