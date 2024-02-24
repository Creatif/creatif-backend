package publish

import "C"
import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/published"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"gorm.io/gorm"
)

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("publish", "Validating...")
	if errs := c.model.Validate(); errs != nil {
		return appErrors.NewValidationError(errs)
	}

	c.logBuilder.Add("publish", "Validated")
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

func (c Main) Logic() (published.Version, error) {
	version := published.NewVersion(c.model.ProjectID)
	if transactionError := storage.Transaction(func(tx *gorm.DB) error {
		if res := tx.Create(&version); res.Error != nil {
			return res.Error
		}

		if res := tx.Exec(getMergeSql(version.ID, (published.PublishedList{}).TableName(), getSelectListSql()), c.model.ProjectID); res.Error != nil {
			return res.Error
		}

		if res := tx.Exec(getMergeSql(version.ID, (published.PublishedMap{}).TableName(), getSelectMapSql()), c.model.ProjectID); res.Error != nil {
			return res.Error
		}

		if res := tx.Exec(getReferenceMergeSql(version.ID, getReferencesSql()), c.model.ProjectID); res.Error != nil {
			return res.Error
		}

		return nil
	}); transactionError != nil {
		return published.Version{}, appErrors.NewApplicationError(transactionError)
	}

	return version, nil
}

func (c Main) Handle() (View, error) {
	if err := c.Validate(); err != nil {
		return View{}, err
	}

	if err := c.Authenticate(); err != nil {
		return View{}, err
	}

	if err := c.Authorize(); err != nil {
		return View{}, err
	}

	model, err := c.Logic()

	if err != nil {
		return View{}, err
	}

	return newView(model), nil
}

func New(model Model, auth auth.Authentication, logBuilder logger.LogBuilder) pkg.Job[Model, View, published.Version] {
	logBuilder.Add("publish", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
