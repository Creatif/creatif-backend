package queryBuilder

import (
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"gorm.io/gorm"
)

type QueryBuilder interface {
	OrderBy() string
	AddWhere(interface{}, ...interface{}) QueryBuilder
	OrderDirection() string
	Limit() int
	Offset() int
	AddSearch(fieldName string, value string) QueryBuilder
	AddIN(fieldName string, value []string) QueryBuilder
	Fields(fields ...string) QueryBuilder
	AddJoin(joinRelation string) QueryBuilder
	AddPreload(preloadRelation string) QueryBuilder
	Build() *gorm.DB
	Run(model interface{}, count *int64, db *gorm.DB) appErrors.AppError[struct{}]
}

type queryBuilder struct {
	orderBy        string
	table          string
	orderDirection string
	limit          int
	offset         int

	search   map[string]string
	in       map[string][]string
	fields   []string
	joins    []string
	preloads []string
	where    []map[interface{}][]interface{}
}

func NewQueryBuilder(table, orderBy, orderDirection string, limit, page int) QueryBuilder {
	if orderBy == "" {
		orderBy = "created_at"
	}

	if orderDirection == "" {
		orderDirection = "asc"
	}

	return &queryBuilder{
		orderBy:        orderBy,
		table:          table,
		orderDirection: orderDirection,
		joins:          make([]string, 0),
		preloads:       make([]string, 0),
		search:         make(map[string]string),
		in:             make(map[string][]string),
		where:          make([]map[interface{}][]interface{}, 0),
		limit:          limit,
		offset:         (page - 1) * limit,
	}
}

func (q *queryBuilder) OrderBy() string {
	return q.orderBy
}

func (q *queryBuilder) OrderDirection() string {
	return q.orderDirection
}

func (q *queryBuilder) Limit() int {
	return q.limit
}

func (q *queryBuilder) Offset() int {
	return q.offset
}

func (q *queryBuilder) Fields(fields ...string) QueryBuilder {
	q.fields = fields

	return q
}

func (q *queryBuilder) AddSearch(fieldName string, value string) QueryBuilder {
	if _, ok := q.search[fieldName]; ok {
		return q
	}

	q.search[fieldName] = value

	return q
}

func (q *queryBuilder) AddIN(fieldName string, value []string) QueryBuilder {
	if _, ok := q.in[fieldName]; ok {
		return q
	}

	q.in[fieldName] = value

	return q
}

func (q *queryBuilder) AddJoin(joinRelation string) QueryBuilder {
	q.joins = append(q.joins, joinRelation)

	return q
}

func (q *queryBuilder) AddWhere(query interface{}, args ...interface{}) QueryBuilder {
	q.where = append(q.where, map[interface{}][]interface{}{
		query: args,
	})

	return q
}

func (q *queryBuilder) AddPreload(preloadRelation string) QueryBuilder {
	q.preloads = append(q.preloads, preloadRelation)

	return q
}

func (q *queryBuilder) Build() *gorm.DB {
	scopes := make([]func(db *gorm.DB) *gorm.DB, 0)
	scopes = append(scopes, storage.WithLimit(q.limit))
	scopes = append(scopes, storage.WithOffset(q.offset))
	scopes = append(scopes, storage.OrderBy(q.orderBy, q.orderDirection))

	if len(q.search) != 0 {
		for fieldName, value := range q.search {
			scopes = append(scopes, storage.WithRegex(fieldName, value))
		}
	}

	if len(q.in) != 0 {
		for fieldName, values := range q.in {
			scopes = append(scopes, storage.WithIN(fieldName, values))
		}
	}

	g := storage.Gorm().
		Table(q.table).
		Scopes(scopes...).
		Select(q.fields)

	if len(q.where) != 0 {
		for _, whereList := range q.where {
			for key, value := range whereList {
				g.Where(key, value...)
			}
		}
	}

	if len(q.joins) != 0 {
		for _, value := range q.joins {
			g.Joins(value)
		}
	}

	if len(q.preloads) != 0 {
		for _, value := range q.preloads {
			g.Preload(value)
		}
	}

	if len(q.fields) != 0 {
		g.Select(q.fields)
	}

	return g
}

func (q *queryBuilder) Run(model interface{}, count *int64, db *gorm.DB) appErrors.AppError[struct{}] {
	if err := runCount(count, q); err != nil {
		return err
	}

	if res := db.Find(model); res.Error != nil {
		return appErrors.NewDatabaseError(res.Error)
	}

	return nil
}

func runCount(count *int64, q *queryBuilder) appErrors.AppError[struct{}] {
	scopes := make([]func(db *gorm.DB) *gorm.DB, 0)

	if len(q.search) != 0 {
		for fieldName, value := range q.search {
			scopes = append(scopes, storage.WithRegex(fieldName, value))
		}
	}

	if len(q.in) != 0 {
		for fieldName, values := range q.in {
			scopes = append(scopes, storage.WithIN(fieldName, values))
		}
	}

	g := storage.Gorm().
		Table(q.table).
		Scopes(scopes...)

	if len(q.where) != 0 {
		for _, whereList := range q.where {
			for key, value := range whereList {
				g.Where(key, value...)
			}
		}
	}

	if res := g.Count(count); res.Error != nil {
		return appErrors.NewDatabaseError(res.Error)
	}

	return nil
}
