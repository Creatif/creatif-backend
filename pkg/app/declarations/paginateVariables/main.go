package paginateVariables

import (
	"creatif/pkg/app/declarations/paginateVariables/pagination"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
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
	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (LogicModel, error) {
	tableName := (declarations.Variable{}).TableName()
	p := pagination.NewPagination(
		tableName,
		fmt.Sprintf("SELECT id, name, behaviour, groups FROM %s", tableName),
		pagination.NewOrderByRule(c.model.Field, c.model.OrderBy, "groups", c.model.Groups),
		c.model.NextID,
		c.model.PrevID,
		c.model.Direction,
		c.model.Limit,
	)
	
	var variables []Variable
	err := p.Paginate(&variables)
	if err != nil {
		return LogicModel{}, appErrors.NewDatabaseError(err).AddError("paginateVariables.declarationsVariable", nil)
	}

	var paginationInfo pagination.PaginationInfo
	if len(variables) == 0 {
		info, err := p.PaginationInfo("", "", c.model.Field, c.model.OrderBy, c.model.Groups, c.model.Limit)
		if err != nil {
			return LogicModel{}, appErrors.NewDatabaseError(err).AddError("paginateVariables.declarationsVariable", nil)
		}

		paginationInfo = info
	} else if len(variables) < c.model.Limit {
		info, err := p.PaginationInfo("", variables[0].ID, c.model.Field, c.model.OrderBy, c.model.Groups, c.model.Limit)
		if err != nil {
			return LogicModel{}, appErrors.NewDatabaseError(err).AddError("paginateVariables.declarationsVariable", nil)
		}

		paginationInfo = info
	} else if len(variables) > 0 {
		info, err := p.PaginationInfo(variables[len(variables)-1].ID, variables[0].ID, c.model.Field, c.model.OrderBy, c.model.Groups, c.model.Limit)
		if err != nil {
			return LogicModel{}, appErrors.NewDatabaseError(err).AddError("paginateVariables.declarationsVariable", nil)
		}

		paginationInfo = info
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
