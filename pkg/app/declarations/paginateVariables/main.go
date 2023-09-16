package paginateVariables

import (
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/sdk/pagination"
	"fmt"
)

type Main struct {
	model Model
}

func (c Main) Validate() error {
	return nil
}
func (c Main) Authenticate() error {
	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (interface{}, error) {
	tableName := (declarations.Variable{}).TableName()
	p := pagination.NewPagination(
		tableName,
		fmt.Sprintf("SELECT id, name, behaviour, groups, metadata FROM %s", tableName),
		pagination.NewOrderByRule(c.model.Field, c.model.OrderBy, "groups", c.model.Groups),
		c.model.NextID,
		c.model.PrevID,
		c.model.Direction,
		c.model.Limit,
	)

	var variables []VariableWithoutValue
	err := p.Paginate(&variables)
	if err != nil {
		return nil, appErrors.NewDatabaseError(err).AddError("paginateVariables.declarationsVariable", nil)
	}

	var paginationInfo pagination.PaginationInfo
	if len(variables) == 0 {
		info, err := p.PaginationInfo("", "")
		if err != nil {
			return nil, appErrors.NewDatabaseError(err).AddError("paginateVariables.declarationsVariable", nil)
		}

		paginationInfo = info
	} else if len(variables) < c.model.Limit {
		info, err := p.PaginationInfo("", variables[0].ID)
		if err != nil {
			return nil, appErrors.NewDatabaseError(err).AddError("paginateVariables.declarationsVariable", nil)
		}

		paginationInfo = info
	} else if len(variables) > 0 {
		info, err := p.PaginationInfo(variables[len(variables)-1].ID, variables[0].ID)
		if err != nil {
			return nil, appErrors.NewDatabaseError(err).AddError("paginateVariables.declarationsVariable", nil)
		}

		paginationInfo = info
	}

	return LogicModelWithoutValue{
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

func New(model Model) pkg.Job[Model, PaginatedView, interface{}] {
	return Main{model: model}
}
