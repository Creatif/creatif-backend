package paginateProjects

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/app"
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

func (c Main) Logic() (sdk.LogicView[QueryModel], error) {
	offset := (c.model.Page - 1) * c.model.Limit
	placeholders := make(map[string]interface{})
	placeholders["offset"] = offset
	placeholders["limit"] = c.model.Limit
	placeholders["user"] = c.auth.User().ID

	countPlaceholders := make(map[string]interface{})
	countPlaceholders["user"] = c.auth.User().ID

	if c.model.OrderBy == "" {
		c.model.OrderBy = "created_at"
	}

	if c.model.OrderDirection == "" {
		c.model.OrderDirection = "ASC"
	}

	c.model.OrderDirection = strings.ToUpper(c.model.OrderDirection)

	var search string
	if c.model.Search != "" {
		search = fmt.Sprintf("AND (%s ILIKE @searchOne OR %s ILIKE @searchTwo OR %s ILIKE @searchThree OR %s ILIKE @searchFour)", "name", "name", "name", "name")
		placeholders["searchOne"] = fmt.Sprintf("%%%s", c.model.Search)
		placeholders["searchTwo"] = fmt.Sprintf("%s%%", c.model.Search)
		placeholders["searchThree"] = fmt.Sprintf("%%%s%%", c.model.Search)
		placeholders["searchFour"] = c.model.Search

		countPlaceholders["searchOne"] = fmt.Sprintf("%%%s", c.model.Search)
		countPlaceholders["searchTwo"] = fmt.Sprintf("%s%%", c.model.Search)
		countPlaceholders["searchThree"] = fmt.Sprintf("%%%s%%", c.model.Search)
		countPlaceholders["searchFour"] = c.model.Search
	}

	sql := fmt.Sprintf(`
SELECT 
    id, 
    name, 
    state,
    created_at,
    updated_at,
    (SELECT count(v.id) FROM declarations.maps AS v WHERE v.project_id = p.id) AS maps_num,
    (SELECT count(v.id) FROM declarations.lists AS v WHERE v.project_id = p.id) AS lists_num
FROM %s AS p
WHERE p.user_id = @user %s
ORDER BY %s %s
OFFSET @offset LIMIT @limit
`, (app.Project{}).TableName(), search, c.model.OrderBy, c.model.OrderDirection)

	var items []QueryModel
	res := storage.Gorm().Raw(sql, placeholders).Scan(&items)
	if res.Error != nil {
		return sdk.LogicView[QueryModel]{}, appErrors.NewDatabaseError(res.Error).AddError("Projects.Paginate.Logic", nil)
	}

	var count int64
	countSql := fmt.Sprintf(`
SELECT count(v.id) AS count
FROM %s AS v
WHERE v.user_id = @user
%s
`, (app.Project{}).TableName(), search)
	res = storage.Gorm().Raw(countSql, countPlaceholders).Scan(&count)
	if res.Error != nil {
		return sdk.LogicView[QueryModel]{}, appErrors.NewDatabaseError(res.Error).AddError("Projects.Paginate.Logic", nil)
	}

	return sdk.LogicView[QueryModel]{
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
		Data:  newView(model.Data),
	}, nil
}

func New(model Model, auth auth.Authentication) pkg.Job[Model, sdk.PaginationView[View], sdk.LogicView[QueryModel]] {
	return Main{model: model, auth: auth}
}
