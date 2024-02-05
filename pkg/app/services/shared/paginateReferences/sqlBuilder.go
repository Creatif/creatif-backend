package paginateReferences

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"fmt"
	"strings"
)

func createSql(model Model, tables [2]string, orderBy, direction, relationshipType string) string {
	returnableFields := ""
	if len(model.Fields) != 0 {
		returnableFields = strings.Join(model.Fields, ",") + ","
	}

	var behaviour string
	if model.Behaviour != "" {
		behaviour = fmt.Sprintf("AND lv.behaviour = @behaviour")
	}

	var locale string
	if len(model.Locales) != 0 {
		resolvedLocales := sdk.Map(model.Locales, func(idx int, value string) string {
			return fmt.Sprintf("'%s'", value)
		})
		locale = fmt.Sprintf("AND lv.locale_id IN(%s)", strings.Join(resolvedLocales, ","))
	}

	var groupsWhereClause string
	if len(model.Groups) != 0 {
		groupsWhereClause = fmt.Sprintf("AND '{%s}'::text[] && %s", strings.Join(model.Groups, ","), "lv.groups")
	}

	var search string
	if model.Search != "" {
		search = fmt.Sprintf("AND (%s ILIKE @searchOne OR %s ILIKE @searchTwo OR %s ILIKE @searchThree OR %s ILIKE @searchFour)", "lv.name", "lv.name", "lv.name", "lv.name")
	}

	relationshipSql := make(map[string]string)
	if relationshipType == "child" {
		relationshipSql["innerJoinOne"] = "r.child_id = lv.id AND r.parent_id = @parentReference AND"
		relationshipSql["innerJoinTwo"] = "AND l.id = @childStructureID AND r.child_structure_id = l.id"
	} else if relationshipType == "parent" {
		relationshipSql["innerJoinOne"] = "r.parent_id = lv.id AND r.parent_id = @parentReference AND r.child_id = @childReference AND"
		relationshipSql["innerJoinTwo"] = "AND l.id = @parentStructureID AND r.parent_structure_id = l.id"
	}

	sql := fmt.Sprintf(`SELECT 
    	lv.id, 
    	lv.short_id, 
    	lv.locale_id,
    	lv.index,
    	lv.name, 
    	lv.behaviour,
    	%s
    	lv.created_at,
    	lv.updated_at
			FROM %s AS r
			INNER JOIN %s AS lv ON %s r.project_id = @projectID
			INNER JOIN %s AS l ON l.project_id = @projectID %s 
		%s %s %s %s
		ORDER BY lv.%s %s
		OFFSET @offset LIMIT @limit`,
		returnableFields,
		(declarations.Reference{}).TableName(),
		tables[0],
		relationshipSql["innerJoinOne"],
		tables[1],
		relationshipSql["innerJoinTwo"],
		locale,
		search,
		groupsWhereClause,
		behaviour,
		orderBy,
		direction)

	return sql
}

func createCountSql(model Model, tables [2]string, relationshipType string) string {
	var behaviour string
	if model.Behaviour != "" {
		behaviour = fmt.Sprintf("AND lv.behaviour = @behaviour")
	}

	var locale string
	if len(model.Locales) != 0 {
		resolvedLocales := sdk.Map(model.Locales, func(idx int, value string) string {
			return fmt.Sprintf("'%s'", value)
		})
		locale = fmt.Sprintf("AND lv.locale_id IN(%s)", strings.Join(resolvedLocales, ","))
	}

	var groupsWhereClause string
	if len(model.Groups) != 0 {
		groupsWhereClause = fmt.Sprintf("AND '{%s}'::text[] && %s", strings.Join(model.Groups, ","), "lv.groups")
	}

	var search string
	if model.Search != "" {
		search = fmt.Sprintf("AND (%s ILIKE @searchOne OR %s ILIKE @searchTwo OR %s ILIKE @searchThree OR %s ILIKE @searchFour)", "lv.name", "lv.name", "lv.name", "lv.name")
	}

	relationshipSql := make(map[string]string)
	if relationshipType == "parent" {
		relationshipSql["innerJoinOne"] = "r.child_id = lv.id AND r.parent_id = @parentReference AND"
		relationshipSql["innerJoinTwo"] = "AND l.id = @childStructureID AND r.child_structure_id = l.id"
	} else if relationshipType == "child" {
		relationshipSql["innerJoinOne"] = "r.parent_id = lv.id AND r.parent_id = @parentReference AND r.child_id = @childReference AND"
		relationshipSql["innerJoinTwo"] = "AND l.id = @parentStructureID AND r.parent_structure_id = l.id"
	}

	sql := fmt.Sprintf(`
    	SELECT 
    	    count(lv.id) AS count
			FROM %s AS r
			INNER JOIN %s AS lv ON %s r.project_id = @projectID
			INNER JOIN %s AS l ON l.project_id = @projectID %s 
		%s %s %s %s
	`,
		(declarations.Reference{}).TableName(),
		tables[0],
		relationshipSql["innerJoinOne"],
		tables[1],
		relationshipSql["innerJoinTwo"],
		locale,
		search,
		behaviour,
		groupsWhereClause,
	)

	return sql
}
