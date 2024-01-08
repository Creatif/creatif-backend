package deleteList

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/shared"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("deleteList", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("deleteList", "Validated")
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
		c.logBuilder.Add("deleteList", res.Error.Error())
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError(res.Error).AddError("deleteList.Logic", nil)
		}

		return nil, appErrors.NewDatabaseError(res.Error).AddError("deleteList.Logic", nil)
	}

	if res.RowsAffected == 0 {
		c.logBuilder.Add("deleteList", "No rows found. That means 404")
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

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, *struct{}, *struct{}] {
	logBuilder.Add("deleteList", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
