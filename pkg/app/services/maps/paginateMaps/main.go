package paginateMaps

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

func (c Main) Logic() (sdk.LogicView[declarations.Map], error) {
	placeholders := make(map[string]interface{})
	placeholders["projectID"] = c.model.ProjectID

	countPlaceholders := make(map[string]interface{})
	countPlaceholders["projectID"] = c.model.ProjectID

	if c.model.OrderBy == "" {
		c.model.OrderBy = "created_at"
	}

	if c.model.OrderDirection == "" {
		c.model.OrderDirection = "ASC"
	}

	c.model.OrderDirection = strings.ToUpper(c.model.OrderDirection)

	var search string
	if c.model.Search != "" {
		search = fmt.Sprintf("AND (%s ILIKE @searchOne OR %s ILIKE @searchTwo OR %s ILIKE @searchThree OR %s ILIKE @searchFour)", "l.name", "l.name", "l.name", "l.name")
		placeholders["searchOne"] = fmt.Sprintf("%%%s", c.model.Search)
		placeholders["searchTwo"] = fmt.Sprintf("%s%%", c.model.Search)
		placeholders["searchThree"] = fmt.Sprintf("%%%s%%", c.model.Search)
		placeholders["searchFour"] = c.model.Search

		countPlaceholders["searchOne"] = fmt.Sprintf("%%%s", c.model.Search)
		countPlaceholders["searchTwo"] = fmt.Sprintf("%s%%", c.model.Search)
		countPlaceholders["searchThree"] = fmt.Sprintf("%%%s%%", c.model.Search)
		countPlaceholders["searchFour"] = c.model.Search
	}

	sql := fmt.Sprintf(`SELECT 
    	l.id, 
    	l.short_id, 
    	l.name, 
    	l.created_at, 
    	l.updated_at 
		FROM %s AS l
		WHERE l.project_id = @projectID
		%s
		ORDER BY l.%s %s`,
		(declarations.Map{}).TableName(),
		search,
		c.model.OrderBy,
		c.model.OrderDirection)

	var items []declarations.Map
	res := storage.Gorm().Raw(sql, placeholders).Scan(&items)
	if res.Error != nil {
		return sdk.LogicView[declarations.Map]{}, appErrors.NewDatabaseError(res.Error).AddError("Maps.Paginate.Logic", nil)
	}

	countSql := fmt.Sprintf(`
    	SELECT 
    	    count(l.id) AS count
		FROM %s AS l
		WHERE l.project_id = @projectID
    	%s
	`,
		(declarations.Map{}).TableName(),
		search,
	)

	var count int64
	res = storage.Gorm().Raw(countSql, countPlaceholders).Scan(&count)
	if res.Error != nil {
		return sdk.LogicView[declarations.Map]{}, appErrors.NewDatabaseError(res.Error).AddError("Maps.Paginate.Logic", nil)
	}

	return sdk.LogicView[declarations.Map]{
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
		return sdk.PaginationView[View]{}, appErrors.NewApplicationError(err).AddError("Maps.Paginate.Handle", nil)
	}

	return sdk.PaginationView[View]{
		Total: model.Total,
		Page:  c.model.Page,
		Data:  items,
	}, nil
}

func New(model Model, auth auth.Authentication) pkg.Job[Model, sdk.PaginationView[View], sdk.LogicView[declarations.Map]] {
	return Main{model: model, auth: auth}
}
