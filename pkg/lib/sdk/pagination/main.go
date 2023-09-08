package pagination

import (
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
)

const ASC = "ASC"
const DESC = "DESC"

const DIRECTION_FORWARD = "forward"
const DIRECTION_BACKWARDS = "backwards"

type SortableIDModel struct {
	ID   string `gorm:"primarykey;type:text CHECK(length(id)=27)"`
	Name string
}

type orderByRule struct {
	field   string
	orderBy string
}

func NewOrderByRule(field, orderBy string) orderByRule {
	return orderByRule{
		field:   field,
		orderBy: orderBy,
	}
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
	isFirstPage := p.cursor == ""
	for i := 0; i < 10; i++ {
		uid, _ := sdk.NewULID()
		fmt.Println(uid)
	}

	if isFirstPage {
		var model SortableIDModel
		if res := storage.Gorm().Raw(fmt.Sprintf("SELECT id, name FROM %s ORDER BY id asc LIMIT 1", p.table)).Scan(&model); res.Error != nil {
			fmt.Println(res.Error)
			return "", res.Error
		}

		fmt.Println(model)
	}

	return "", nil
}
