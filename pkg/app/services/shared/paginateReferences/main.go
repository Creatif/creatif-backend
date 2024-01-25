package paginateReferences

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"strings"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
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

func (c Main) Logic() (sdk.LogicView[declarations.MapVariable], error) {
	offset := (c.model.Page - 1) * c.model.Limit
	placeholders := make(map[string]interface{})
	placeholders["projectID"] = c.model.ProjectID
	placeholders["offset"] = offset
	placeholders["limit"] = c.model.Limit
	placeholders["parentReference"] = c.model.ParentID

	countPlaceholders := make(map[string]interface{})
	countPlaceholders["projectID"] = c.model.ProjectID

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
		groupsWhereClause = fmt.Sprintf("WHERE '{%s}'::text[] && %s", strings.Join(c.model.Groups, ","), "lv.groups")
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
	if len(c.model.Fields) != 0 {
		returnableFields = strings.Join(c.model.Fields, ",") + ","
	}

	tabels := []string{
		(declarations.MapVariable{}).TableName(),
		(declarations.Map{}).TableName(),
	}

	if c.model.StructureType == "list" {
		tabels = []string{
			(declarations.MapVariable{}).TableName(),
			(declarations.Map{}).TableName(),
		}
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
			INNER JOIN %s AS lv
			INNER JOIN %s AS l
		ON l.project_id = @projectID AND l.id = lv.map_id AND r.parent_id = lv.id AND r.parent_id = @parentReference %s %s
		%s
		%s
		ORDER BY lv.%s %s
		OFFSET @offset LIMIT @limit`,
		returnableFields,
		(declarations.Reference{}).TableName(),
		tabels[0],
		tabels[1],
		locale,
		search,
		groupsWhereClause,
		behaviour,
		c.model.OrderBy,
		c.model.OrderDirection)

	var items []declarations.MapVariable
	res := storage.Gorm().Raw(sql, placeholders).Scan(&items)
	if res.Error != nil {
		c.logBuilder.Add("paginateMapVariables", res.Error.Error())
		return sdk.LogicView[declarations.MapVariable]{}, appErrors.NewDatabaseError(res.Error).AddError("Maps.Paginate.Logic", nil)
	}

	return sdk.LogicView[declarations.MapVariable]{
		Total: 0,
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, sdk.PaginationView[View], sdk.LogicView[declarations.MapVariable]] {
	logBuilder.Add("paginateMapVariables", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
