package removeMap

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
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

func (c Main) Logic() (interface{}, error) {
	res := storage.Gorm().Where(fmt.Sprintf("project_id = ? AND (name = ? OR id = ? OR short_id = ?)"), c.model.ProjectID, c.model.Name, c.model.Name, c.model.Name).Delete(&declarations.Map{})
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError(res.Error).AddError("removeMap.Logic", nil)
		}

		return nil, appErrors.NewDatabaseError(res.Error).AddError("removeMap.Logic", nil)
	}

	if res.RowsAffected == 0 {
		return nil, appErrors.NewNotFoundError(errors.New(fmt.Sprintf("Map with name '%s' not found.", c.model.Name))).AddError("removeMap.Logic", nil)
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, interface{}, interface{}] {
	return Main{model: model, auth: auth}
}
