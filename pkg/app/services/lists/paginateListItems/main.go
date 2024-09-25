package paginateListItems

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"strings"
)

type Main struct {
	model Model
	auth  auth.Authentication
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	return nil
}

func (c Main) Authenticate() error {
	if err := c.auth.Authenticate(); err != nil {
		return appErrors.NewAuthenticationError(err)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (sdk.LogicView[QueryVariable], error) {
	offset := (c.model.Page - 1) * c.model.Limit
	placeholders := make(map[string]interface{})
	placeholders["projectID"] = c.model.ProjectID
	placeholders["offset"] = offset
	placeholders["name"] = c.model.ListName
	placeholders["limit"] = c.model.Limit
	placeholders["groups"] = c.model.Groups

	countPlaceholders := make(map[string]interface{})
	countPlaceholders["projectID"] = c.model.ProjectID
	countPlaceholders["name"] = c.model.ListName
	countPlaceholders["groups"] = c.model.Groups

	if c.model.OrderBy == "" {
		c.model.OrderBy = "index"
	}

	var behaviour string
	if c.model.Behaviour != "" {
		behaviour = fmt.Sprintf("AND lv.behaviour = @behaviour")
		placeholders["behaviour"] = c.model.Behaviour
		countPlaceholders["behaviour"] = c.model.Behaviour
	}

	var locale string
	if len(c.model.Locales) != 0 {
		resolvedLocales := sdk.Map(c.model.Locales, func(idx int, value string) string {
			return fmt.Sprintf("'%s'", value)
		})
		locale = fmt.Sprintf("AND lv.locale_id IN(%s)", strings.Join(resolvedLocales, ","))
	}

	if c.model.OrderDirection == "" {
		c.model.OrderDirection = "ASC"
	}

	c.model.OrderDirection = strings.ToUpper(c.model.OrderDirection)

	var groupsWhereClause string
	if len(c.model.Groups) != 0 {
		searchForGroups := strings.Join(c.model.Groups, ",")
		groupsWhereClause = fmt.Sprintf("INNER JOIN LATERAL (SELECT g.variable_id, g.group_id, g.groups FROM %s AS g WHERE lv.id = g.variable_id ORDER BY g.variable_id LIMIT 1) AS g ON '{%s}'::text[] && g.groups", (declarations.VariableGroup{}).TableName(), searchForGroups)
	}

	var search string
	if c.model.Search != "" {
		search = fmt.Sprintf("AND (%s ILIKE @searchOne OR %s ILIKE @searchTwo OR %s ILIKE @searchThree OR %s ILIKE @searchFour)", "lv.name", "lv.name", "lv.name", "lv.name")
		placeholders["searchOne"] = fmt.Sprintf("%%%s", c.model.Search)
		placeholders["searchTwo"] = fmt.Sprintf("%s%%", c.model.Search)
		placeholders["searchThree"] = fmt.Sprintf("%%%s%%", c.model.Search)
		placeholders["searchFour"] = c.model.Search

		countPlaceholders["searchOne"] = fmt.Sprintf("%%%s", c.model.Search)
		countPlaceholders["searchTwo"] = fmt.Sprintf("%s%%", c.model.Search)
		countPlaceholders["searchThree"] = fmt.Sprintf("%%%s%%", c.model.Search)
		countPlaceholders["searchFour"] = c.model.Search
	}

	returnableFields := ""
	groupsSubquery := ""
	if len(c.model.Fields) != 0 {
		if sdk.Includes(c.model.Fields, "groups") {
			groupsSubquery = fmt.Sprintf("ARRAY((SELECT g.name FROM declarations.groups AS g INNER JOIN declarations.variable_groups AS vg ON vg.group_id = g.id AND vg.variable_id = lv.id)) AS groups")
		}

		returnableFields = strings.Join(sdk.Filter(c.model.Fields, func(idx int, value string) bool {
			return value != "groups"
		}), ",") + ","
	}

	sql := fmt.Sprintf(`SELECT 
    	lv.id, 
    	lv.index, 
    	lv.short_id, 
    	lv.locale_id,
    	lv.name, 
    	lv.behaviour, 
    	%s
    	%s
    	lv.created_at, 
    	lv.updated_at 
			FROM %s AS lv
			INNER JOIN %s AS l
		ON l.project_id = @projectID AND (l.id = @name OR l.short_id = @name) AND l.id = lv.list_id %s %s
		%s
		%s
		ORDER BY lv.%s %s
		OFFSET @offset LIMIT @limit`,
		groupsSubquery,
		returnableFields,
		(declarations.ListVariable{}).TableName(),
		(declarations.List{}).TableName(),
		locale,
		search,
		groupsWhereClause,
		behaviour,
		c.model.OrderBy,
		c.model.OrderDirection)

	var items []QueryVariable
	res := storage.Gorm().Raw(sql, placeholders).Scan(&items)
	if res.Error != nil {
		return sdk.LogicView[QueryVariable]{}, appErrors.NewDatabaseError(res.Error).AddError("ListItems.Paginate.Logic", nil)
	}

	countSql := fmt.Sprintf(`
    	SELECT 
    	    count(lv.id) AS count
		FROM %s AS lv
		INNER JOIN %s AS l
		ON l.project_id = @projectID AND (l.short_id = @name OR l.id = @name) AND l.id = lv.list_id %s %s
    	%s
    	%s
	`,
		(declarations.ListVariable{}).TableName(),
		(declarations.List{}).TableName(),
		locale,
		search,
		behaviour,
		groupsWhereClause,
	)

	var count int64
	res = storage.Gorm().Raw(countSql, countPlaceholders).Scan(&count)
	if res.Error != nil {
		return sdk.LogicView[QueryVariable]{}, appErrors.NewDatabaseError(res.Error).AddError("ListItems.Paginate.Logic", nil)
	}

	return sdk.LogicView[QueryVariable]{
		Total: count,
		Data:  items,
	}, nil
}

func (c Main) Handle() (sdk.PaginationView[View], error) {
	if err := c.Validate(); err != nil {
		return sdk.PaginationView[View]{}, err
	}

	if err := c.Authenticate(); err != nil {
		return sdk.PaginationView[View]{}, err
	}

	if err := c.Authorize(); err != nil {
		return sdk.PaginationView[View]{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return sdk.PaginationView[View]{}, err
	}

	items, err := newView(model.Data)
	if err != nil {
		return sdk.PaginationView[View]{}, appErrors.NewApplicationError(err).AddError("ListItems.Paginate.Handle", nil)
	}

	return sdk.PaginationView[View]{
		Total: model.Total,
		Page:  c.model.Page,
		Data:  items,
	}, nil
}

func New(model Model, auth auth.Authentication) pkg.Job[Model, sdk.PaginationView[View], sdk.LogicView[QueryVariable]] {
	return Main{model: model, auth: auth}
}
