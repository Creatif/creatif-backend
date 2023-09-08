package pagination

import (
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/sdk/pagination"
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
	p := pagination.NewPagination((declarations.Node{}).TableName(), "nothing", pagination.NewOrderByRule(c.model.SortField, c.model.SortOrder), "")
	p.Create()

	return nil, nil
}

func (c Main) Handle() ([]View, error) {
	if err := c.Validate(); err != nil {
		return []View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return []View{}, err
	}

	if err := c.Authorize(); err != nil {
		return []View{}, err
	}

	_, err := c.Logic()

	if err != nil {
		return []View{}, err
	}

	return []View{}, nil
}

func New(model PaginationModel) pkg.Job[PaginationModel, []View, interface{}] {
	return Main{model: model}
}
