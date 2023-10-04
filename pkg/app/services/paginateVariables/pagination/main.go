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
	projectId    string
	table        string
	sql          string
	paginationId string
	limit        int
	rule         orderByRule
	direction    string
}

func NewPagination(projectId, table, sql string, rules orderByRule, paginationId, direction string, limit int) *Pagination {
	return &Pagination{
		projectId:    projectId,
		table:        table,
		sql:          sql,
		rule:         rules,
		paginationId: paginationId,
		direction:    direction,
		limit:        limit,
	}
}

func (p Pagination) Paginate(model interface{}) error {
	isFirstPage := p.paginationId == ""
	if isFirstPage {
		id, err := getInitialID(p.projectId, p.table, p.rule.orderBy)
		if err != nil {
			return err
		}

		operator, err := getOperator(DIRECTION_FORWARD, p.rule.orderBy, true)
		if err != nil {
			return err
		}

		if len(p.rule.groups) > 0 {
			groups := strings.Join(p.rule.groups, ",")
			if res := storage.Gorm().Raw(fmt.Sprintf("%s WHERE project_id = '%s' AND id %s '%s' AND '{%s}'::text[] && %s ORDER BY %s %s LIMIT %d", p.sql, p.projectId, operator, id, groups, p.rule.groupField, p.rule.field, p.rule.orderBy, p.limit)).Scan(model); res.Error != nil {
				return res.Error
			}
		} else {
			sql := fmt.Sprintf("%s WHERE project_id = '%s' AND id %s '%s' ORDER BY %s %s LIMIT %d", p.sql, p.projectId, operator, id, p.rule.field, p.rule.orderBy, p.limit)
			if res := storage.Gorm().Raw(sql).Scan(model); res.Error != nil {
				return res.Error
			}
		}

		return nil
	} else {
		operator, err := getOperator(p.direction, p.rule.orderBy, false)
		if err != nil {
			return err
		}

		sql := fmt.Sprintf("%s WHERE project_id = '%s' AND id %s '%s' ORDER BY %s %s LIMIT %d", p.sql, p.projectId, operator, p.paginationId, p.rule.field, p.rule.orderBy, p.limit)
		if res := storage.Gorm().Raw(sql).Scan(model); res.Error != nil {
			return res.Error
		}
	}

	return nil
}

func (p Pagination) PaginationInfo(prevPaginationId, paginationId, field, orderBy string, groups []string, limit int) (PaginationInfo, error) {
	var next, prev string
	if paginationId != "" {
		next = fmt.Sprintf("?paginationId=%s&field=%s&orderBy=%s&direction=%s&limit=%d", paginationId, p.rule.field, p.rule.orderBy, DIRECTION_FORWARD, p.limit)
	}

	if prevPaginationId != "" {
		next = fmt.Sprintf("?paginationId=%s&field=%s&orderBy=%s&direction=%s&limit=%d", paginationId, p.rule.field, p.rule.orderBy, DIRECTION_BACKWARDS, p.limit)
	}

	if len(groups) == 0 {
		groups = make([]string, 0)
	}

	return PaginationInfo{
		Next: next,
		Prev: prev,
		Parameters: Parameters{
			PaginationID: paginationId,
			Field:        field,
			OrderBy:      orderBy,
			Groups:       groups,
			Limit:        limit,
		},
	}, nil
}
