package deleteListItemByID

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/locales"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"gorm.io/gorm"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
}

func (c Main) Validate() error {
	c.logBuilder.Add("deleteListItemByID", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("deleteListItemByID", "Validated")
	return nil
}

func (c Main) Authenticate() error {
	// user check by project id should be gotten here, with authentication cookie
	var project app.Project
	if err := storage.Get((app.Project{}).TableName(), c.model.ProjectID, &project); err != nil {
		return appErrors.NewAuthenticationError(err).AddError("createVariable.Authenticate", nil)
	}

	return nil
}

func (c Main) Authorize() error {
	return nil
}

func (c Main) Logic() (*struct{}, error) {
	localeID, err := locales.GetIDWithAlpha(c.model.Locale)
	if err != nil {
		c.logBuilder.Add("deleteListItemByID", err.Error())
		return nil, appErrors.NewApplicationError(err).AddError("deleteListItemByID.Logic", nil)
	}
	var list declarations.List
	res := storage.Gorm().Where("name = ? AND project_id = ? AND locale_id = ?", c.model.Name, c.model.ProjectID, localeID).Select("ID").First(&list)
	if res.Error != nil {
		c.logBuilder.Add("deleteListItemByID", res.Error.Error())
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError(res.Error).AddError("deleteListItemByID.Logic", nil)
		}

		return nil, appErrors.NewDatabaseError(res.Error).AddError("deleteListItemByID.Logic", nil)
	}

	var variable declarations.ListVariable
	res = storage.Gorm().Where("id = ? AND list_id = ? AND locale_id = ?", c.model.ItemID, list.ID, localeID).Delete(&variable)
	if res.Error != nil {
		c.logBuilder.Add("deleteListItemByID", res.Error.Error())

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFoundError(res.Error).AddError("deleteListItemByID.Logic", nil)
		}

		return nil, appErrors.NewDatabaseError(res.Error).AddError("deleteListItemByID.Logic", nil)
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

func New(model Model, logBuilder logger.LogBuilder) pkg.Job[Model, *struct{}, *struct{}] {
	logBuilder.Add("deleteListItemByID", "Created")
	return Main{model: model, logBuilder: logBuilder}
}
