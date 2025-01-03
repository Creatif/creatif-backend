package deleteList

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/shared"
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

func (c Main) Logic() (*struct{}, error) {
	id, val := shared.DetermineID("", c.model.Name, c.model.ID, c.model.ShortID)
	var list declarations.List
	res := storage.Gorm().Where(fmt.Sprintf("%s AND project_id = ?", id), val, c.model.ProjectID).Delete(&list)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError(res.Error).AddError("deleteList.Logic", nil)
		}

		return nil, appErrors.NewDatabaseError(res.Error).AddError("deleteList.Logic", nil)
	}

	if err := removeConnections(c.model.ProjectID, c.model.ID); err != nil {
		return nil, appErrors.NewNotFoundError(res.Error).AddError("deleteList.Logic", nil)
	}

	if res.RowsAffected == 0 {
		return nil, appErrors.NewNotFoundError(res.Error).AddError("deleteList.Logic", nil)
	}

	return nil, nil
}

func (c Main) Handle() (*struct{}, error) {
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

func New(model Model, auth auth.Authentication) pkg.Job[Model, *struct{}, *struct{}] {
	return Main{model: model, auth: auth}
}
