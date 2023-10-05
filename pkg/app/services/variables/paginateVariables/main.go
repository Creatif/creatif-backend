package paginateVariables

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	pagination2 "creatif/pkg/app/services/variables/paginateVariables/pagination"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/sdk"
	"creatif/pkg/lib/storage"
	"fmt"
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

func (c Main) Logic() (LogicModel, error) {
	tableName := (declarations.Variable{}).TableName()
	p := pagination2.NewPagination(
		c.model.ProjectID,
		tableName,
		fmt.Sprintf("SELECT id, name, behaviour, groups FROM %s", tableName),
		pagination2.NewOrderByRule(c.model.Field, c.model.OrderBy, "groups", c.model.Groups),
		c.model.PaginationID,
		c.model.Direction,
		c.model.Limit,
	)

	var variables []Variable
	err := p.Paginate(&variables)
	if err != nil {
		return LogicModel{}, appErrors.NewDatabaseError(err).AddError("paginateVariables.Logic", nil)
	}

	paginationId, err := pagination2.ResolveCursor(c.model.PaginationID, c.model.Direction, c.model.OrderBy, sdk.Map(variables, func(idx int, value Variable) string {
		return value.ID
	}), c.model.Limit)

	if err != nil {
		return LogicModel{}, appErrors.NewApplicationError(err).AddError("paginateVariables.Logic", nil)
	}

	paginationInfo, err := p.PaginationInfo(c.model.PaginationID, paginationId, c.model.Field, c.model.OrderBy, c.model.Groups, c.model.Limit)
	if err != nil {
		return LogicModel{}, appErrors.NewDatabaseError(err).AddError("paginateVariables.Logic", nil)
	}

	return LogicModel{
		variables:      variables,
		paginationInfo: paginationInfo,
	}, nil
}

func (c Main) Handle() (PaginatedView, error) {
	if err := c.Validate(); err != nil {
		return PaginatedView{}, err
	}

	if err := c.Authenticate(); err != nil {
		return PaginatedView{}, err
	}

	if err := c.Authorize(); err != nil {
		return PaginatedView{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return PaginatedView{}, err
	}

	return newView(model), nil
}

func New(model Model) pkg.Job[Model, PaginatedView, LogicModel] {
	return Main{model: model}
}