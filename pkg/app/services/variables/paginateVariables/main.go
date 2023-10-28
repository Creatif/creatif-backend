package paginateVariables

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
	c.logBuilder.Add("paginateVariables", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}
	c.logBuilder.Add("paginateVariables", "Validated")
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

func (c Main) Logic() (sdk.LogicView[declarations.Variable], error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("paginateVariables", err.Error())
		return sdk.LogicView[declarations.Variable]{}, appErrors.NewApplicationError(err).AddError("Variables.Paginate.Logic", nil)
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
		groupsWhereClause = fmt.Sprintf("AND '{%s}'::text[] && %s", strings.Join(c.model.Groups, ","), "groups")
	}

	sql := fmt.Sprintf(`
SELECT id, short_id, groups, name, behaviour, metadata, value, created_at, updated_at
FROM declarations.variables
WHERE locale_id = ? AND project_id = ?
%s
ORDER BY %s %s
OFFSET ? LIMIT ?
`, groupsWhereClause, c.model.OrderBy, c.model.OrderDirection)

	offset := (c.model.Page - 1) * c.model.Limit
	var items []declarations.Variable
	res := storage.Gorm().Raw(sql, localeID, c.model.ProjectID, offset, c.model.Limit).Scan(&items)
	if res.Error != nil {
		c.logBuilder.Add("paginateVariables", res.Error.Error())
		return sdk.LogicView[declarations.Variable]{}, appErrors.NewDatabaseError(err).AddError("Variables.Paginate.Logic", nil)
	}

	var count int64
	countSql := fmt.Sprintf(`
SELECT count(v.id) AS count
FROM declarations.variables AS v
WHERE v.locale_id = ? AND v.project_id = ?
%s
`, groupsWhereClause)
	res = storage.Gorm().Raw(countSql, localeID, c.model.ProjectID).Scan(&count)
	if res.Error != nil {
		c.logBuilder.Add("paginateVariables", res.Error.Error())
		return sdk.LogicView[declarations.Variable]{}, appErrors.NewDatabaseError(err).AddError("Variables.Paginate.Logic", nil)
	}

	return sdk.LogicView[declarations.Variable]{
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

	return sdk.PaginationView[View]{
		Total: model.Total,
		Page:  c.model.Page,
		Data:  newView(model.Data, c.model.Locale),
	}, nil
}

func New(model Model, logBuilder logger.LogBuilder) pkg.Job[Model, sdk.PaginationView[View], sdk.LogicView[declarations.Variable]] {
	logBuilder.Add("paginateVariables", "Created.")
	return Main{model: model, logBuilder: logBuilder}
}
