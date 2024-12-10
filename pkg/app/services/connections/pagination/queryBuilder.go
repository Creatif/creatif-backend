package pagination

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"strings"
	"time"
)

type defaults struct {
	orderBy         string
	orderDirections string
}

type subQueries struct {
	behaviour        string
	locale           string
	groups           string
	search           string
	returnableFields string
	groupsSubQuery   string
}

type QueryVariable struct {
	ID      string  `gorm:"primarykey;type:text"`
	ShortID string  `gorm:"uniqueIndex:unique_map_variable;type:text;not null"`
	Index   float64 `gorm:"type:float"`

	Name      string         `gorm:"uniqueIndex:unique_map_variable;not null"`
	Behaviour string         `gorm:"not null"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	Value     datatypes.JSON `gorm:"type:jsonb"`
	Groups    pq.StringArray `gorm:"type:text[];column:groups"`

	MapID    string `gorm:"uniqueIndex:unique_map_variable;type:text"`
	LocaleID string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"<-:create;index"`
	UpdatedAt time.Time
}

func createQueryPlaceholders(projectId, structureId, parentVariableId string, offset int, limit int, groups []string, behaviour, search string) map[string]interface{} {
	placeholders := make(map[string]interface{})

	placeholders["projectID"] = projectId
	placeholders["offset"] = offset
	placeholders["structureId"] = structureId
	placeholders["parentVariableId"] = parentVariableId
	placeholders["limit"] = limit
	placeholders["groups"] = groups

	if behaviour != "" {
		placeholders["behaviour"] = behaviour
	}

	if search != "" {
		placeholders["searchOne"] = fmt.Sprintf("%%%s", search)
		placeholders["searchTwo"] = fmt.Sprintf("%s%%", search)
		placeholders["searchThree"] = fmt.Sprintf("%%%s%%", search)
		placeholders["searchFour"] = search
	}

	return placeholders
}

func createCountPlaceholders(projectId, structureId, parentVariableId string, groups []string, behaviour, search string) map[string]interface{} {
	placeholders := make(map[string]interface{})

	placeholders["projectID"] = projectId
	placeholders["structureId"] = structureId
	placeholders["parentVariableId"] = parentVariableId
	placeholders["groups"] = groups

	if behaviour != "" {
		placeholders["behaviour"] = behaviour
	}

	if search != "" {
		placeholders["searchOne"] = fmt.Sprintf("%%%s", search)
		placeholders["searchTwo"] = fmt.Sprintf("%s%%", search)
		placeholders["searchThree"] = fmt.Sprintf("%%%s%%", search)
		placeholders["searchFour"] = search
	}

	return placeholders
}

func createDefaults(orderBy, orderDirection string) defaults {
	var def defaults
	def.orderBy = "index"

	if orderBy != "" {
		def.orderBy = orderBy
	}

	if orderDirection == "" {
		orderDirection = "ASC"
	}

	def.orderDirections = strings.ToUpper(orderDirection)

	return def
}

func createSubQueries(behaviour string, locales, groups []string, search string) subQueries {
	var sq subQueries

	if behaviour != "" {
		sq.behaviour = fmt.Sprintf("AND lv.behaviour = @behaviour")
	}

	if len(locales) != 0 {
		resolvedLocales := sdk.Map(locales, func(idx int, value string) string {
			return fmt.Sprintf("'%s'", value)
		})
		sq.locale = fmt.Sprintf("AND lv.locale_id IN(%s)", strings.Join(resolvedLocales, ","))
	}

	if len(groups) != 0 {
		searchForGroups := strings.Join(groups, ",")
		sq.groups = fmt.Sprintf("INNER JOIN %s AS g ON lv.id = g.variable_id AND '{%s}'::text[] && g.groups", (declarations.VariableGroup{}).TableName(), searchForGroups)
	}

	if search != "" {
		sq.search = fmt.Sprintf("AND (%s ILIKE @searchOne OR %s ILIKE @searchTwo OR %s ILIKE @searchThree OR %s ILIKE @searchFour)", "lv.name", "lv.name", "lv.name", "lv.name")
	}

	return sq
}

func createPaginationSql(structureType string, sq subQueries, defs defaults) string {
	sql := fmt.Sprintf(`SELECT 
    	lv.id, 
    	lv.index, 
    	lv.short_id, 
    	lv.locale_id,
    	lv.name, 
    	lv.behaviour,
    	%s
    	lv.created_at, 
    	lv.updated_at 
			FROM %s AS lv
			INNER JOIN %s AS l ON l.project_id = @projectID AND l.id = @structureId AND l.id = lv.list_id 
			INNER JOIN %s AS c ON c.parent_variable_id = @parentVariableId AND c.child_variable_id = lv.id
		%s 
		%s
		%s
		%s
		ORDER BY lv.%s %s
		OFFSET @offset LIMIT @limit`,
		sq.returnableFields,
		(declarations.ListVariable{}).TableName(),
		(declarations.List{}).TableName(),
		(declarations.Connection{}).TableName(),
		sq.locale,
		sq.search,
		sq.groups,
		sq.behaviour,
		defs.orderBy,
		defs.orderDirections,
	)

	if structureType == "map" {
		sql = fmt.Sprintf(`SELECT 
    	lv.id, 
    	lv.index, 
    	lv.short_id, 
    	lv.locale_id,
    	lv.name, 
    	lv.behaviour,
    	%s
    	lv.created_at, 
    	lv.updated_at 
			FROM %s AS lv
			INNER JOIN %s AS l ON l.project_id = @projectID AND l.id = @structureId AND l.id = lv.map_id 
			INNER JOIN %s AS c ON c.parent_variable_id = @parentVariableId AND c.child_variable_id = lv.id
		%s 
		%s
		%s
		%s
		ORDER BY lv.%s %s
		OFFSET @offset LIMIT @limit`,
			sq.returnableFields,
			(declarations.MapVariable{}).TableName(),
			(declarations.Map{}).TableName(),
			(declarations.Connection{}).TableName(),
			sq.locale,
			sq.search,
			sq.groups,
			sq.behaviour,
			defs.orderBy,
			defs.orderDirections,
		)
	}

	return sql
}

func createCountSql(structureType string, sq subQueries) string {
	sql := fmt.Sprintf(`
    	SELECT 
    	    count(lv.id) AS count
		FROM %s AS lv
		INNER JOIN %s AS l ON l.project_id = @projectID AND l.id = @structureId AND l.id = lv.list_id 
		INNER JOIN %s AS c ON c.parent_variable_id = @parentVariableId AND c.child_variable_id = lv.id
		%s 
		%s
    	%s
    	%s
	`,
		(declarations.ListVariable{}).TableName(),
		(declarations.List{}).TableName(),
		(declarations.Connection{}).TableName(),
		sq.locale,
		sq.search,
		sq.behaviour,
		sq.groups,
	)

	if structureType == "map" {
		sql = fmt.Sprintf(`
    	SELECT 
    	    count(lv.id) AS count
		FROM %s AS lv
		INNER JOIN %s AS l ON l.project_id = @projectID AND l.id = @structureId AND l.id = lv.map_id
		INNER JOIN %s AS c ON c.parent_variable_id = @parentVariableId AND c.child_variable_id = lv.id
		%s 
		%s
    	%s
    	%s
	`,
			(declarations.MapVariable{}).TableName(),
			(declarations.Map{}).TableName(),
			(declarations.Connection{}).TableName(),
			sq.locale,
			sq.search,
			sq.behaviour,
			sq.groups,
		)
	}

	return sql
}
