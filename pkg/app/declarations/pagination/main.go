package pagination

import (
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/sdk/pagination"
	"fmt"
)

type Main struct {
	model PaginationModel
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
	tableName := (declarations.Node{}).TableName()
	p := pagination.NewPagination(
		tableName,
		fmt.Sprintf("SELECT id, name, behaviour, groups, metadata FROM %s", tableName),
		pagination.NewOrderByRule(c.model.Field, c.model.OrderBy, "groups", c.model.Groups),
		c.model.NextID,
		c.model.PrevID,
		c.model.Direction,
		c.model.Limit,
	)

	var nodes []NodeWithoutValue
	err := p.Paginate(&nodes)
	if err != nil {
		return nil, appErrors.NewDatabaseError(err).AddError("pagination.declarationsNode", nil)
	}

	var paginationInfo pagination.PaginationInfo
	if len(nodes) == 0 {
		info, err := p.PaginationInfo("", "")
		if err != nil {
			return nil, appErrors.NewDatabaseError(err).AddError("pagination.declarationsNode", nil)
		}

		paginationInfo = info
	} else if len(nodes) < c.model.Limit {
		info, err := p.PaginationInfo("", nodes[0].ID)
		if err != nil {
			return nil, appErrors.NewDatabaseError(err).AddError("pagination.declarationsNode", nil)
		}

		paginationInfo = info
	} else if len(nodes) > 0 {
		info, err := p.PaginationInfo(nodes[len(nodes)-1].ID, nodes[0].ID)
		if err != nil {
			return nil, appErrors.NewDatabaseError(err).AddError("pagination.declarationsNode", nil)
		}

		paginationInfo = info
	}

	return LogicModelWithoutValue{
		nodes:          nodes,
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

func New(model PaginationModel) pkg.Job[PaginationModel, PaginatedView, interface{}] {
	return Main{model: model}
}
