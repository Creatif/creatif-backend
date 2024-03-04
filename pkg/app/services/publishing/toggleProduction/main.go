package toggleProduction

import (
	"creatif/pkg/app/auth"
	"creatif/pkg/app/domain/published"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Results struct {
	Errors []error
}

type Main struct {
	model      Model
	logBuilder logger.LogBuilder
	auth       auth.Authentication
}

func (c Main) Validate() error {
	c.logBuilder.Add("removeVersion", "Validating...")
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

	var version published.Version
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT id, is_production_version FROM %s WHERE project_id = ? AND id = ?", (published.Version{}).TableName()), c.model.ProjectID, c.model.ID).Scan(&version)

	if res.Error != nil {
		return nil, appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected == 0 {
		return nil, appErrors.NewNotFoundError(errors.New("This version does not exist."))
	}

	if version.IsProductionVersion {
		if res := storage.Gorm().Exec(fmt.Sprintf("UPDATE %s SET is_production_version = false WHERE project_id = ? AND id = ?", (published.Version{}).TableName()), c.model.ProjectID, c.model.ID); res.Error != nil {
			return nil, appErrors.NewApplicationError(res.Error)
		}

		return nil, nil
	}

	if transactionError := storage.Transaction(func(tx *gorm.DB) error {
		if res := storage.Gorm().Exec(fmt.Sprintf("UPDATE %s SET is_production_version = false WHERE project_id = ? AND is_production_version = true", (published.Version{}).TableName()), c.model.ProjectID); res.Error != nil {
			return res.Error
		}

		if res := storage.Gorm().Exec(fmt.Sprintf("UPDATE %s SET is_production_version = true WHERE project_id = ? AND id = ?", (published.Version{}).TableName()), c.model.ProjectID, c.model.ID); res.Error != nil {
			return res.Error
		}

		return nil
	}); transactionError != nil {
		return nil, appErrors.NewApplicationError(transactionError)
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
	logBuilder.Add("removeVersion", "Created")
	return Main{model: model, logBuilder: logBuilder, auth: auth}
}
