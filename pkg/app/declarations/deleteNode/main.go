package deleteNode

import (
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
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
	if res := storage.Gorm().Where("name = ?", c.model.Name).Delete(&declarations.Node{}); res.Error != nil {
		return nil, appErrors.NewDatabaseError(res.Error).AddError("deleteNode.Logic", nil)
	}
	return nil, nil
}

func (c Main) Handle() (interface{}, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if err := c.Authenticate(); err != nil {
		return nil, err
	}

	if err := c.Authorize(); err != nil {
		return nil, err
	}

	_, err := c.Logic()

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func New(model Model) pkg.Job[Model, interface{}, interface{}] {
	return Main{model: model}
}
