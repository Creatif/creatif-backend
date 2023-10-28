package paginateListItems

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
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
}

func (c Main) Validate() error {
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	return nil
}
func (c Main) Authenticate() error {
	// user check by project id should be gotten here, with authentication cookie
	var project app.Project
	if err := storage.Get((app.Project{}).TableName(), c.model.ProjectID, &project); err != nil {
		return appErrors.NewAuthenticationError(err).AddError("paginateVariables.Authenticate", nil)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (sdk.LogicView[declarations.ListVariable], error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("paginateListItems", err.Error())
		return sdk.LogicView[declarations.ListVariable]{}, appErrors.NewApplicationError(err).AddError("ListItems.Paginate.Logic", nil)
	}

	if c.model.OrderBy == "" {
		c.model.OrderBy = "created_at"
	}

	if c.model.OrderDirection == "" {
		c.model.OrderDirection = "ASC"
	}

	c.model.OrderDirection = strings.ToUpper(c.model.OrderDirection)

	var groupsWhereClause string
	if len(c.model.Groups) != 0 {
		groupsWhereClause = fmt.Sprintf("WHERE '{%s}'::text[] && %s", strings.Join(c.model.Groups, ","), "lv.groups")
	}

	offset := (c.model.Page - 1) * c.model.Limit
	sql := fmt.Sprintf(`SELECT 
    	lv.id, 
    	lv.index, 
    	lv.short_id, 
    	lv.locale_id,
    	lv.name, 
    	lv.groups,
    	lv.behaviour, 
    	lv.metadata, 
    	lv.value, 
    	lv.created_at, 
    	lv.updated_at 
			FROM %s AS lv
			INNER JOIN %s AS l
		ON l.project_id = ? AND l.name = ? AND l.id = lv.list_id AND l.locale_id = ?
		%s
		ORDER BY lv.%s %s
		OFFSET ? LIMIT ?`,
		(declarations.ListVariable{}).TableName(),
		(declarations.List{}).TableName(),
		groupsWhereClause,
		c.model.OrderBy,
		c.model.OrderDirection)

	var items []declarations.ListVariable
	res := storage.Gorm().Raw(sql, c.model.ProjectID, c.model.ListName, localeID, offset, c.model.Limit).Scan(&items)
	if res.Error != nil {
		c.logBuilder.Add("paginateListItems", res.Error.Error())
		return sdk.LogicView[declarations.ListVariable]{}, appErrors.NewDatabaseError(res.Error).AddError("ListItems.Paginate.Logic", nil)
	}

	countSql := fmt.Sprintf(`
    	SELECT 
    	    count(lv.id) AS count
		FROM %s AS lv
		INNER JOIN %s AS l
		ON l.project_id = ? AND l.name = ? AND l.id = lv.list_id
    	%s
	`,
		(declarations.ListVariable{}).TableName(),
		(declarations.List{}).TableName(),
		groupsWhereClause,
	)

	var count int64
	res = storage.Gorm().Raw(countSql, c.model.ProjectID, c.model.ListName).Scan(&count)
	if res.Error != nil {
		c.logBuilder.Add("paginateListItem", res.Error.Error())
		return sdk.LogicView[declarations.ListVariable]{}, appErrors.NewDatabaseError(res.Error).AddError("ListItems.Paginate.Logic", nil)
	}

	return sdk.LogicView[declarations.ListVariable]{
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

func New(model Model, logBuilder logger.LogBuilder) pkg.Job[Model, sdk.PaginationView[View], sdk.LogicView[declarations.ListVariable]] {
	logBuilder.Add("paginateListItems", "Created")
	return Main{model: model, logBuilder: logBuilder}
}
