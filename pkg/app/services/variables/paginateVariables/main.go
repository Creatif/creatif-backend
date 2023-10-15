package paginateVariables

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/queryBuilder"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
	"strings"
)

type Main struct {
	model Model
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

func (c Main) Logic() (sdk.LogicView[declarations.Variable], error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		return sdk.LogicView[declarations.Variable]{}, appErrors.NewApplicationError(err).AddError("Variables.Paginate.Logic", nil)
	}

	qb := queryBuilder.NewQueryBuilder(fmt.Sprintf("%s as v", (declarations.Variable{}).TableName()), c.model.OrderBy, c.model.OrderDirection, c.model.Limit, c.model.Page)
	qb = qb.Fields("v.id", "v.short_id", "v.groups", "v.name", "v.behaviour", "v.metadata", "v.value", "v.created_at", "v.updated_at").
		AddWhere("v.locale_id = ?", localeID).
		AddWhere("v.project_id = ?", c.model.ProjectID)

	if len(c.model.Groups) != 0 {
		qb.AddWhere(fmt.Sprintf("'{%s}'::text[] && %s", strings.Join(c.model.Groups, ","), "groups"))
	}

	var items []declarations.Variable
	var count int64
	if err := qb.Run(&items, &count, qb.Build()); err != nil {
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

func New(model Model) pkg.Job[Model, sdk.PaginationView[View], sdk.LogicView[declarations.Variable]] {
	return Main{model: model}
}
